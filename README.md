# demon-logger
Connects to a TCP service and logs the result into file. [Download](https://github.com/yene/demon-logger/releases/latest)


Example, run for 5 days and write to file every 5 minutes.

`./demon-logger-osx -host 127.0.0.1 -flush 300 -age 5`

### Help
```
Usage of ./demon-logger-osx:
  -age int
        How long the app runs in days. (default 2)
  -flush int
        Write to disk interval in seconds (default 300)
  -host string
        Host IP (default "127.0.0.1")
```


### Build for ARM
`env GOOS=linux GOARCH=arm GOARM=7 go build -o "demon-logger-arm"`

### keep alive
It can not be determined without a ping/keepalive if the network is still working. As solution use the keepalive package or SetReadDeadline.
I am not sure if the keepalive package works perfectly, but it is good enough.

### References
* [Reconnect TCP](http://stackoverflow.com/questions/23395519/reconnect-tcp-on-eof-in-go)
* [Debug go routines](http://stackoverflow.com/a/19145992/279890)
* [tcp keepalive](https://github.com/felixge/tcpkeepalive)
* [How to kill children](http://stackoverflow.com/a/6807784/279890)

### License
MIT
