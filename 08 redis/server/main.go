package main

import (
	"fmt"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Received request from", conn.RemoteAddr())
	conn.Write([]byte("Hello, World!"))
}
