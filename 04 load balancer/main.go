package main

import (
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()

	fmt.Println("1Server is listening on port 8080...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		fmt.Println("n", n, err)
		if err != nil {
			fmt.Println("Connection closed.")
			return
		}
		if n == 0 {
			fmt.Println("Connection closed.")
			conn.Write([]byte("Connection closed."))
			conn.Close()
			return
		}
		fmt.Println("message", string(buf[:n]))
		conn.Write([]byte("Message received."))
		return
	}
}
