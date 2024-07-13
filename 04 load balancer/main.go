package main

import (
	"bufio"
	"bytes"
	"fmt"
	"lb/serversmanager"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

func main() {
	timeout := time.Second * 5
	conf := []serversmanager.Config{
		{
			Url:     "localhost:8080",
			Timeout: timeout,
		},
		{
			Url:     "localhost:8081",
			Timeout: timeout,
		},
		{
			Url:     "localhost:8082",
			Timeout: timeout,
		},
	}
	servers := make([]*serversmanager.ServerManager, 0)
	for _, c := range conf {
		servers = append(servers, serversmanager.NewServerManager(c))
	}
	serversPool := serversmanager.NewServerPool(
		servers, 3)

	ln, err := net.Listen("tcp", ":4000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("LB Listening on localhost:4000")
	defer ln.Close()
	serversPool.BootupServers()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handler(conn, serversPool)
	}
}

func handler(conn net.Conn, serversPool *serversmanager.ServerPool) {
	fmt.Println("Received request from", conn.RemoteAddr())
	defer conn.Close()

	// read from the connection
	clientRes, err := readFromConn(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Client request is:\n%s\n", clientRes)
	server := serversPool.GetNextServer()
	if server == nil {
		fmt.Println("No server available")
		// wrtie a response to the client
		buf := bytes.Buffer{}
		buf.WriteString("HTTP/1.1 502 Bad Gateway\r\n")
		buf.WriteString("\r\n")
		buf.WriteString("All connnections are dead")
		conn.Write(buf.Bytes())
		conn.Close()
		return
	}

	backConn, err := server.Dial(3)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer backConn.Close()
	// forward the request to the server
	backConn.Write([]byte(clientRes))
	// read the response from the server
	// and forward it to the client
	backRes, err := readFromConn(backConn)
	fmt.Printf("backRes: \n%s\n", backRes)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write([]byte(backRes))

}

func readFromConn(conn net.Conn) (string, error) {
	reader := bufio.NewReader(conn)
	buf := bytes.Buffer{}
	contentLength := 0
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			break
		}
		buf.WriteString(s)
		if strings.HasPrefix(s, "Content-Length:") {
			len := strings.TrimPrefix(s, "Content-Length:")
			len = strings.TrimSpace(len)
			contentLength, err = strconv.Atoi(len)
			if err != nil {
				log.Fatal(err)
			}
		}
		if s == "\r\n" {
			break
		}
	}
	for contentLength > 0 {
		b, err := reader.ReadByte()
		if err != nil {
			fmt.Println(err)
			break
		}
		buf.WriteByte(b)
		contentLength--
	}
	return buf.String(), nil
}
