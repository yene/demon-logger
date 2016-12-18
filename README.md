# demon-logger
A small app that connects to a TCP service and logs the result into file.

`./demon-logger -host 192.168.1.47`

### keep alive
It can not be determined without a ping/keepalive if the network is still working. As solution use the keepalive package or SetReadDeadline.
I am not sure if the keepalive package works perfectly, but it is good enough.

### References
* [Reconnect TCP](http://stackoverflow.com/questions/23395519/reconnect-tcp-on-eof-in-go)
* [Debug go routines](http://stackoverflow.com/a/19145992/279890)
* [tcp keepalive](https://github.com/felixge/tcpkeepalive)
* [How to kill children](http://stackoverflow.com/a/6807784/279890)
