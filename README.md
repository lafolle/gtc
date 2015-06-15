
###Go Timed Command

`gtc` runs a command for a specificed period of time.
If no time is specified, command is run until it exits.

###Install
`go get github.com/lafolle/gtc`

###Usage
```
Usage of gtc:
  -p=false: perserve status of cmd
  -t=2522880h: duration of cmd
```

###Todo
1. Tests
2. Improvise on the design of communication betweeen main and gocoroutine
3. Relay signals to `gtc` to the command being executed

###References:
1. http://blog.golang.org/go-concurrency-patterns-timing-out-and
2. https://golang.org/pkg/os/exec/
3. https://golang.org/pkh/os
