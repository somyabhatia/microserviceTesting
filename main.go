package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

type config struct {
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

	flag.StringVar(&laddr, "l", "", "listening address")
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
		go recv(conn, c)
	}
}

func recv(conn net.Conn, c config) {
	addr := remoteAddr(conn)
	log.Infof("accepted connection from %s", addr)
	r := bufio.NewReader(conn)
	for {
		conn.SetDeadline(time.Now().Add(c.timeout))
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Errorln(err)
			break
		}
		if _, err := conn.Write([]byte("\n")); err != nil {
			log.Errorln(err)
			break
		}
		fmt.Printf("%s --> %s", addr, msg)
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
		addr := remoteAddr(conn)
		log.Infof("connected to %s", addr)
		r := bufio.NewReader(conn)
		for {
			conn.SetDeadline(time.Now().Add(c.timeout))
			msg := fmt.Sprintf("%d\n", i)
			if _, err := conn.Write([]byte(msg)); err != nil {
				log.Errorln(err)
				break
			}
			_, err := r.ReadString('\n')
			if err != nil {
				log.Errorln(err)
				break
			}
			fmt.Printf("%s <-- %s", addr, msg)
			i++
			time.Sleep(c.sendInterval)
		}
		conn.Close()
	}
}

func remoteAddr(conn net.Conn) string {
	ip := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	names, err := net.LookupAddr(ip)
	if err == nil && len(names) > 0 {
		parts := strings.Split(names[0], ".")
		return parts[0] + " (" + ip + ")"
	}
	return ip
}
