
###Go Timed Command

`gtc` runs the command for the specificed period of time.
If no time is specified, the command is run until it exits.(or until 290 years)

###Install
`go get github.com/lafolle/gtc`

###Usage
```
Usage of gtc:
  -p=false: perserve status of cmd
  -t=9.152s: duration of cmd
```

###Todo
1. Tests

###References:
1. http://blog.golang.org/go-concurrency-patterns-timing-out-and
2. https://golang.org/pkg/os/exec/
3. https://golang.org/pkh/os
