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
		listen bool
		laddr  string
		c      config
	)

	flag.BoolVar(&listen, "l", false, "listen")
	flag.StringVar(&laddr, "a", "", "listening address")
	flag.DurationVar(&c.timeout, "t", 5*time.Second, "connect/send/recv timeout")
	flag.DurationVar(&c.retryInterval, "r", 1*time.Second, "connection retry interval")
	flag.DurationVar(&c.sendInterval, "i", 1*time.Second, "message sending interval")
	flag.Parse()
	addrs := flag.Args()

	if listen {
		laddr = normalizeAddr(laddr)
		log.Infof("listening on %s", laddr)
		ln, err := net.Listen("tcp", laddr)
		if err != nil {
			log.Fatal(err)
		}
		go accept(ln, c)
	}

	for _, addr := range addrs {
		go send(normalizeAddr(addr), c)
	}

	select {}
}

func normalizeAddr(addr string) string {
	parts := strings.Split(addr, ":")
	if len(parts) == 1 {
		return addr + ":8080"
	}
	return addr
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
	log.Infof("accepted connection from %s", remoteAddr(conn))
	r := bufio.NewReader(conn)
	for {
		conn.SetDeadline(time.Now().Add(c.timeout))
		s, err := r.ReadString('\n')
		parts := strings.Split(s, "\n")
		msg := parts[0]
		if err != nil {
			log.Errorln(err)
			break
		}
		if _, err := conn.Write([]byte("\n")); err != nil {
			log.Errorln(err)
			break
		}
		fmt.Printf("%s <-- %s\n", msg, remoteAddr(conn))
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
		log.Infof("connected to %s", remoteAddr(conn))
		r := bufio.NewReader(conn)
		for {
			conn.SetDeadline(time.Now().Add(c.timeout))
			msg := fmt.Sprintf("%d", i)
			if _, err := conn.Write([]byte(msg + "\n")); err != nil {
				log.Errorln(err)
				break
			}
			_, err := r.ReadString('\n')
			if err != nil {
				log.Errorln(err)
				break
			}
			fmt.Printf("%s --> %s\n", msg, remoteAddr(conn))
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
