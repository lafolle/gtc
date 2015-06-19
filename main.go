// gtc runs a command for specified period of time
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

const (
	DEFAULT_CMD_DURATION time.Duration = 2522880 * time.Hour
)

// attachToStdReadPipe channelizes data from `in` to `out`
func attachToStdReadPipe(in io.Reader, out *os.File) {
	r := bufio.NewReader(in)
	for {
		str, err := r.ReadString('\n')
		fmt.Fprint(out, str)
		if err != nil {
			break
		}
	}
}

// handles signals sent to gtc. All signals to gtc are forwarded
// to cmd being executed. If due to the signal cmd exits gtc is
// unable to send signal to cmd, gtc exits. In latter case, with
// an error.
func handleSignals(cmd *exec.Cmd, sigchan chan os.Signal, done chan error) {
	var e error
	for {
		sig := <-sigchan
		if err := cmd.Process.Signal(sig); err != nil {
			e = err
			break
		} else {
			if cmd.ProcessState != nil && cmd.ProcessState.Exited() {
				// if cmd has exited due to a signal, it probably won't come
				// here, because, cmd.Wait() returns imidiately. In case it does:
				e = nil
				break
			}
		}
	}
	done <- e
}

// startcmd get the standard output and error pipes
// and runs the command. It reads from pipes and spits
// output to stdout/stderr of gtc
func startCmd(cmd *exec.Cmd, done chan error, sigchan chan os.Signal) {
	if stdout, err := cmd.StdoutPipe(); err == nil {
		go attachToStdReadPipe(stdout, os.Stdout)
	} else {
		done <- err
	}
	if stderr, err := cmd.StderrPipe(); err == nil {
		go attachToStdReadPipe(stderr, os.Stderr)
	} else {
		done <- err
	}
	if err := cmd.Start(); err != nil {
		done <- err
	}
	go handleSignals(cmd, sigchan, done)

	done <- cmd.Wait()
}

// RunCmd runs command `cmd` for at most `duration` duration
func RunCmd(cmd *exec.Cmd, duration *time.Duration, preservestatus *bool) {
	done := make(chan error, 1)
	sigchan := make(chan os.Signal, 1)
	defer close(done)
	defer close(sigchan)
	signal.Notify(sigchan)

	go startCmd(cmd, done, sigchan)

	select {
	case err := <-done:
		if err != nil {
			log.Fatal(err)
		}
		if *preservestatus == true {
			ws, _ := cmd.ProcessState.Sys().(syscall.WaitStatus)
			os.Exit(ws.ExitStatus())
		}
	case <-time.After(*duration):
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal(err)
		}
	}
}

func showHelp() {
	fmt.Println(`Usage of gtc:
  -p=false: perserve status of cmd
  -t=2522880h0m0s: duration to run cmd
  command`)
}

// set commands checks the flag.Args() and creates a
// command structs out of it.
func SetCmd(c []string) *exec.Cmd {
	var cmd *exec.Cmd
	if len(c) == 0 {
		showHelp()
		return nil
	} else if len(c) == 1 {
		cmd = exec.Command(c[0])
	} else {
		cmd = exec.Command(c[0], c[1:]...)
	}
	return cmd
}

func main() {
	var (
		// default val is no. of hours in 288 years
		dur            *time.Duration = flag.Duration("t", DEFAULT_CMD_DURATION, "duration of cmd") 
		preservestatus *bool          = flag.Bool("p", false, "perserve status of cmd")
	)
	flag.Parse()
	cmd := SetCmd(flag.Args())
	if cmd == nil {
		os.Exit(1)
	}
	RunCmd(cmd, dur, preservestatus)
}

/*
	[1] What program does with a signal is different for different processes.
		Like on TERM signal, one program might sucide, while other might just
		go to background. Simplest way to check is, after sending signal to process
		check if its state is still running. If yes, let it run, else, declare
		completion of process.
*/
