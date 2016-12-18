package main

import (
	"bufio"
	"errors"
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

const port = "1039"
const filename = "log.txt"
const retryInterval = 5

var age = flag.Int("age", 2, "How long the app runs in days.")
var host = flag.String("host", "127.0.0.1", "Host IP")
var interval = flag.Int("flush", 3600, "Flush interval in seconds")

func main() {
	flag.Parse()
	errCh := make(chan error)
	for {
		conn, err := net.Dial("tcp", *host+":"+port)
		if err != nil {
			log.Println("Could not connect to demon", err)
		} else {
			kaConn, _ := tcpkeepalive.EnableKeepAlive(conn)
			kaConn.SetKeepAliveIdle(30 * time.Second)
			kaConn.SetKeepAliveCount(4)
			kaConn.SetKeepAliveInterval(5 * time.Second)
			go readLog(conn, errCh)
			readErr := <-errCh
			log.Println("Error", readErr)
			conn.Close()
		}
		log.Println("retrying in", retryInterval, "seconds")
		time.Sleep(retryInterval * time.Second)
	}
}

func readLog(conn net.Conn, errCh chan<- error) {
	reader := bufio.NewReader(conn)

	f, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
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
		fmt.Println("Logger expired after days:", *age)
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
			errCh <- errors.New("Could not read from the conncetion: " + err.Error())
			return
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
