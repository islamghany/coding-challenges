package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"lb/serversmanager"
	"lb/serversmanager/algorithms"
	"log"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	mux sync.Mutex
)

type config struct {
	algorithm string
	timeout   time.Duration
}

func main() {
	var cfg config
	flag.StringVar(&cfg.algorithm, "algorithm", "round_robin", "Load balancing algorithm")
	flag.DurationVar(&cfg.timeout, "timeout", time.Second*5, "Server timeout")
	flag.Parse()
	// Check if the algorithm is available
	if err := algorithms.IsAlgorithmAvailable(cfg.algorithm); err != nil {
		log.Fatalf("Error: %s", err)
	}
	// Create server pool
	conf := []serversmanager.Config{
		{
			Url:     "localhost:8080",
			Timeout: cfg.timeout,
		},
		{
			Url:     "localhost:8081",
			Timeout: cfg.timeout,
		},
		{
			Url:     "localhost:8082",
			Timeout: cfg.timeout,
		},
	}
	servers := make([]*serversmanager.ServerManager, 0)
	for _, c := range conf {
		servers = append(servers, serversmanager.NewServerManager(c))
	}
	serversPool := serversmanager.NewServerPool(
		servers, 3)
	// Start the load balancer server
	ln, err := net.Listen("tcp", ":4000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("LB Listening on localhost:4000")
	defer ln.Close()
	// Bootup the servers
	serversPool.BootupServers()
	// Load balancer algorithm selection
	var lbalgo algorithms.LoadBalancerAlgorithm
	switch cfg.algorithm {
	case algorithms.ROUNDROBIN:
		lbalgo = algorithms.NewRoundRobin()
	case algorithms.LEASTCONNECTION:
		lbalgo = algorithms.NewLeastConnection()
	default:
		log.Fatalf("Algorithm not implemented")
	}
	// Accept incoming connections
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go handler(conn, serversPool, lbalgo)
	}
}

func handler(conn net.Conn, serversPool *serversmanager.ServerPool, lbalgo algorithms.LoadBalancerAlgorithm) {

	fmt.Println("Received request from", conn.RemoteAddr())
	defer conn.Close()

	// read from the connection
	clientRes, err := readFromConn(conn)
	if err != nil {
		fmt.Println(err)
		return
	}
	servers := serversPool.GetAllActiveServers()
	fmt.Printf("Active servers: %d\n", len(servers))

	if len(servers) == 0 {
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
	server := lbalgo.Next(servers)
	backConn, err := server.Dial(3)
	if err != nil {
		fmt.Println(err)
		return
	}
	// startTime := time.Now()
	server.IncrementTotalRequests()
	defer func() {
		server.DecrementTotalRequests()
		backConn.Close()
	}()
	// forward the request to the server
	backConn.Write([]byte(clientRes))
	// read the response from the server
	// and forward it to the client
	backRes, err := readFromConn(backConn)
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.Write([]byte(backRes))
	fmt.Println("Served by", server.GeURL())

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
