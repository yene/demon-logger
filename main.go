package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/felixge/tcpkeepalive"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

const PORT = "1039"
const FILENAME = "log.txt"

var age = flag.Int("age", 2, "How long the app runs in days.")
var host = flag.String("host", "127.0.0.1", "Host IP")
var interval = flag.Int("flush", 3600, "Flush interval in seconds")

func main() {
	flag.Parse()
	// TODO: loop over do manual retry, block with a channel
	// http://stackoverflow.com/questions/23395519/reconnect-tcp-on-eof-in-go
	errCh := make(chan error)
	for {
		conn, err := net.Dial("tcp", *host+":"+PORT)
		if err != nil {
			log.Println("Could not connect to demon", err)
		} else {
			kaConn, _ := tcpkeepalive.EnableKeepAlive(conn)
			kaConn.SetKeepAliveIdle(30 * time.Second)
			kaConn.SetKeepAliveCount(4)
			kaConn.SetKeepAliveInterval(5 * time.Second)
			readLog(conn, errCh)
			err = <-errCh
			log.Println("Error", err)
			conn.Close()
		}
		log.Println("retrying in 10 seconds")
		time.Sleep(30 * time.Second)
	}
}

func readLog(conn net.Conn, errCh chan error) {
	reader := bufio.NewReader(conn)

	f, err := os.OpenFile(FILENAME, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)

	t := time.Now().Format("2006-01-02 15:04:05")
	fmt.Fprintln(w, t+" Connected to "+*host)

	flushticker := time.NewTicker(time.Second * time.Duration(*interval))
	go func() {
		for t := range flushticker.C {
			fmt.Println("Writting file to disk at", t)
			err := w.Flush()
			check(err)
		}
	}()

	expiretimer := time.NewTimer(time.Hour * 24 * time.Duration(*age))
	go func() {
		<-expiretimer.C
		err := w.Flush()
		check(err)
		println("Logger expired after days:", *age)
		os.Exit(0)
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			err := w.Flush()
			check(err)
			println("Stopping because of interrupt", sig)
			os.Exit(1)
		}
	}()
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Could not connect to demon", err)
		}
		fmt.Fprintln(w, line)
		fmt.Println(line)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
