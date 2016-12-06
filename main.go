package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/textproto"
	"os"
	"os/signal"
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

const PORT = "1039"

var age = flag.Int("age", 2, "How long the app runs in days.")
var host = flag.String("host", "127.0.0.1", "Host IP")
var interval = flag.Int("flush", 3600, "Flush interval in seconds")

func main() {
	flag.Parse()
	// TODO: loop over do manual retry, block with a channel
	// http://stackoverflow.com/questions/23395519/reconnect-tcp-on-eof-in-go
	conn, err := net.Dial("tcp", *host+":"+PORT)
	if err != nil {
		log.Println("Could not connect to demon", err)
		os.Exit(1)
	}
	defer conn.Close()

	var buffer bytes.Buffer

	check(err)

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	f, err := os.OpenFile("file.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)

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
		line, _ := tp.ReadLine()
		fmt.Fprintln(w, line)
		fmt.Println(line)
		buffer.WriteString(line + "\n")
		err := ioutil.WriteFile("test", buffer.Bytes(), 0644)
		check(err)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
