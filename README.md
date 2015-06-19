
###Go Timed Command

`gtc` runs the command for at most given duration.
If no duration is provided, `gtc` runs the command unitl it exits by itself.

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

###Test checks
1. gtc -t=12s ls -R /
	should run the command for 12 secs
2. gtc ls -R /
	should run the command until it exits by itself
3. gtc ls -R / > /tmp/k
	should write to /tmp/k from stdout of cmd
4. echo / | gtc -t=3s ls -R 
	this works bit differently.
	`echo / | xargs gtc -t=3s ls -R `
	This needs to be investigated as currently gtc does not 
	its stdin to `cmd`'s stdin

###References:
1. http://blog.golang.org/go-concurrency-patterns-timing-out-and
2. https://golang.org/pkg/os/exec/
3. https://golang.org/pkh/os
