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
	"time"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

const PORT = "1039"

var duration = flag.Int("duration", 2, "How long the app runs in days.")

func main() {
	flag.Parse()

	expiretimer := time.NewTimer(time.Hour * 24 * time.Duration(*duration))
	go func() {
		<-expiretimer.C
		println("Logger expired after days:", duration)
		os.Exit(0)
	}()

	host := flag.Arg(0)
	conn, err := net.Dial("tcp", host+":"+PORT)
	if err != nil {
		// handle error
		log.Println("Could not connect to demon", err)
		os.Exit(1)
	}
	defer conn.Close()

	var buffer bytes.Buffer

	check(err)

	reader := bufio.NewReader(conn)
	tp := textproto.NewReader(reader)

	f, err := os.Create("file.txt")
	check(err)
	defer f.Close()
	w := bufio.NewWriter(f)

	flushticker := time.NewTicker(time.Hour)
	go func() {
		for t := range flushticker.C {
			fmt.Println("Writting file to disk at", t)
			err := w.Flush()
			check(err)
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
