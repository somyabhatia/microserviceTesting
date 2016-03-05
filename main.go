package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	log "github.com/Sirupsen/logrus"
)

type config struct {
	msg           string
	timeout       time.Duration
	retryInterval time.Duration
	sendInterval  time.Duration
}

func main() {
	log.SetOutput(os.Stderr)
	var (
		laddr string
		c     config
	)
	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	flag.StringVar(&laddr, "l", "", "listening address")
	flag.StringVar(&c.msg, "m", hostname, "message to send")
	flag.DurationVar(&c.timeout, "t", 5*time.Second, "connect/send/recv timeout")
	flag.DurationVar(&c.retryInterval, "r", 1*time.Second, "connection retry interval")
	flag.DurationVar(&c.sendInterval, "i", 1*time.Second, "message sending interval")
	flag.Parse()
	addrs := flag.Args()

	if laddr != "" {
		log.Infof("listening on %s", laddr)
		ln, err := net.Listen("tcp", laddr)
		if err != nil {
			log.Fatal(err)
		}
		go accept(ln, c)
	}

	for _, addr := range addrs {
		go send(addr, c)
	}

	select {}
}

func accept(ln net.Listener, c config) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Errorf("accept failed: %s", err)
			continue
		}
		log.Infof("accepted connection from %s", conn.RemoteAddr())
		go recv(conn, c)
	}
}

func recv(conn net.Conn, c config) {
	r := bufio.NewReader(conn)
	for {
		conn.SetDeadline(time.Now().Add(c.timeout))
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Errorln(err)
			break
		}
		fmt.Print(msg)
		if _, err := conn.Write([]byte("\n")); err != nil {
			log.Errorln(err)
			break
		}
	}
	conn.Close()
}

func send(addr string, c config) {
	i := 1
	for {
		conn, err := net.DialTimeout("tcp", addr, c.timeout)
		if err != nil {
			log.Errorln(err)
			time.Sleep(c.retryInterval)
			continue
		}
		log.Infof("connected to %s", conn.RemoteAddr())
		r := bufio.NewReader(conn)
		for {
			conn.SetDeadline(time.Now().Add(c.timeout))
			if _, err := conn.Write([]byte(fmt.Sprintf("%s %d\n", c.msg, i))); err != nil {
				log.Errorln(err)
				break
			}
			_, err := r.ReadString('\n')
			if err != nil {
				log.Errorln(err)
				break
			}
			i++
			time.Sleep(c.sendInterval)
		}
		conn.Close()
	}
}
