package main

import (
	"fmt"
	"log"
	"crypto/tls"
	"net"
	"bufio"
	"strconv"
	"strings"
)

func main() {
	var webPort int = 443
	log.SetFlags(log.Lshortfile)

	cer, err := tls.LoadX509KeyPair("keys/server.crt", "keys/server.key")
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Keys loaded!")

	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", joinStr(":", strconv.Itoa(webPort)), config)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println(joinStr("Allocated Port ",strconv.Itoa(webPort)))
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("New Connection")
	r := bufio.NewReader(conn)
	for {
		msg, err := r.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		println(msg)

		n, err := conn.Write([]byte("world\n"))
		if err != nil {
			log.Println(n, err)
			return
		}
	}
}

func joinStr(strs ...string) string {
	var sb strings.Builder
	for _, str := range strs {
		sb.WriteString(str)
	}
	return sb.String()
}