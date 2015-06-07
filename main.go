package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"time"
)

var (
	dur            *time.Duration = flag.Duration("t", 9152000000, "duration of cmd")
	preservestatus *bool          = flag.Bool("p", false, "perserve status of cmd")
)

func startcmd(cmd *exec.Cmd, done chan error) {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	outreader := bufio.NewReader(stdout)
	go func(r *bufio.Reader) {
		for {
			str, err := r.ReadString('\n')
			fmt.Print(str)
			if err != nil {
				break
			}
		}
	}(outreader)
	errreader := bufio.NewReader(stderr)
	go func(r *bufio.Reader) {
		for {
			str, err := r.ReadString('\n')
			fmt.Print(str)
			if err != nil {
				break
			}
		}
	}(errreader)
	done <- cmd.Wait()
}

func setcmd(c []string) *exec.Cmd {
	var cmd *exec.Cmd
	if len(c) == 0 {
		fmt.Println("No command provided. Aborting.")
		return nil
	} else if len(c) == 1 {
		cmd = exec.Command(c[0])
	} else {
		cmd = exec.Command(c[0], c[1:]...)
	}
	return cmd
}

func main() {
	flag.Parse()

	cmd := setcmd(flag.Args())
	if cmd == nil {
		os.Exit(1)
	}

	done := make(chan error, 1)
	defer close(done)
	go startcmd(cmd, done)

	select {
	case <-done:
		// TODO: for making gtc cross platform type of sys() need to be checked
		ws, _ := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exit_status := ws.ExitStatus()
		if *preservestatus == true {
			os.Exit(exit_status)
		}
	case <-time.After(*dur):
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal(err)
		}
	}
}
