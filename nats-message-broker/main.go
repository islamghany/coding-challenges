package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"nats/commands"
	"nats/parser"
	"net"
	"time"
)

func main() {

	listner, err := net.Listen("tcp", ":4222")
	if err != nil {
		panic(err)
	}
	for {
		conn, err := listner.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				log.Printf("Temporary error accepting connection %v", netErr)
				continue
			}
			log.Fatalf("Error accepting connection %v", err)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		err := conn.Close()
		fmt.Println("Connection closed")
		if err != nil {
			fmt.Println("Error closing connection", err.Error())
		}
	}()
	reader := bufio.NewReader(conn)
	for {
		// Set a read deadline to prevent hanging connections
		if err := conn.SetReadDeadline(time.Now().Add(5 * time.Minute)); err != nil {
			log.Printf("Error setting read deadline: %v", err)
			return
		}
		cmd, err := parser.Parse(reader)

		if err != nil {
			if err == io.EOF {
				log.Printf("Client closed connection")
				return
			}
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Read timeout")
				return
			}
			log.Printf("Error parsing command: %v", err)
			if _, err := conn.Write([]byte("ERROR\r\n")); err != nil {
				log.Printf("Error writing error response: %v", err)
				return
			}
			continue
		}

		commands.HandleCommand(cmd, conn)
	}
}
