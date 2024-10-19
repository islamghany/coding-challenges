package main

import (
	"fmt"
	"httproxy/http"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	PORT = "8989"
)

var (
	forbiddenHosts = make(map[string]bool)
)

func main() {

	// load forbidden hosts
	forbiddenHosts["facebook.com"] = true

	// Create a socket
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", PORT))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	// Graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		log.Println("Shutting down the server...")
		listener.Close()
		os.Exit(0)
	}()

	// Accept incoming connections
	log.Printf("Starting proxy server on port %s\n", PORT)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(clientConn net.Conn) {
	defer func() {
		log.Println("Closing client connection")
		clientConn.Close()
	}()

	// Read the incoming request
	request, err := http.ParseIncomingRequest(clientConn)
	if err != nil {
		log.Println(err)
		return
	}
	// Print the request
	fmt.Printf("%+v\n", request)

	host := request.GetHeader("Host")
	log.Printf("Request from %s, Target %s\n", clientConn.RemoteAddr(), host)
	if forbiddenHosts[host] {
		res := http.NewResponse(nil)
		res.SetStatus(403)
		res.SetHeader("X-Content-Type-Options", "nosniff") // Security header to prevent MIME type sniffing
		res.SetHeader("Content-Type", "text/plain")
		res.Body = []byte("Access to this host is forbidden")
		http.WriteResponseToConn(clientConn, res)
		return
	}
	// Create a connection to the target host
	targetConn, err := net.Dial("tcp", fmt.Sprintf("%s:80", host))
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		log.Println("Closing target connection")
		targetConn.Close()

	}()

	// Add X-Forwarded-For header
	clientIP := strings.Split(clientConn.RemoteAddr().String(), ":")[0]
	request.SetHeader("X-Forwarded-For", clientIP)

	// Remove hop-by-hop headers
	request.RemoveHeader("Proxy-Connection")

	// Write the request to the target host
	err = http.WriteRequestToConn(targetConn, request)
	if err != nil {
		log.Println("Error writing request:", err)
		return
	}

	// Read the response from the target host
	_, err = io.Copy(clientConn, targetConn)
	if err != nil {
		log.Println("Error copying response:", err)
		return
	}

}
