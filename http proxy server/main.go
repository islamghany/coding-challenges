package main

import (
	"bufio"
	"fmt"
	"httproxy/http"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

const (
	PORT               = "8989"
	forbiddenHostsFile = "forbidden-hosts.txt"
	bannedWordsFile    = "banned-words.txt"
)

var forbiddenHosts map[string]bool
var bannedWords map[string]bool

func main() {

	// Load forbidden hosts
	forbiddenHosts = loadFileLinesAsMap(forbiddenHostsFile)

	// Load banned words
	bannedWords = loadFileLinesAsMap(bannedWordsFile)

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

	// Check if the request is a CONNECT method
	if request.Method == "CONNECT" {
		handleConnectMethod(clientConn, request)
		return
	}

	// Add X-Forwarded-For header
	clientIP := strings.Split(clientConn.RemoteAddr().String(), ":")[0]
	request.SetHeader("X-Forwarded-For", clientIP)

	// Remove hop-by-hop headers
	request.RemoveHeader("Proxy-Connection")

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

	// Write the request to the target host
	err = http.WriteRequestToConn(targetConn, request)
	if err != nil {
		log.Println("Error writing request:", err)
		return
	}

	// Read the response from the target host
	response, err := http.ParseIncomingResponse(targetConn)
	if err != nil {
		log.Println("Error reading response:", err)
		return
	}

	if containsBannedWords(response.Body) {
		// Return a 403 Forbidden response
		forbiddenResponse := http.NewResponse(nil)
		response.SetStatus(403)
		forbiddenResponse.SetHeader("Content-Type", "text/plain; charset=utf-8")
		forbiddenResponse.SetHeader("X-Content-Type-Options", "nosniff")
		forbiddenResponse.Body = []byte("Website content not allowed.")
		http.WriteResponseToConn(clientConn, forbiddenResponse)
		return
	}

	fmt.Printf("%+v\n", response)

	err = http.WriteResponseToConn(clientConn, response)
	if err != nil {
		log.Println("Error writing response:", err)
		return
	}

}

func handleConnectMethod(clientConn net.Conn, request *http.Request) {
	wg := sync.WaitGroup{}
	// Extract the host and port from the CONNECT request
	targetAddress := request.GetHeader("Host")
	if targetAddress == "" {
		log.Println("Invalid CONNECT request, missing Host header")
		return
	}
	// Establish a connection to the target server
	targetConn, err := net.Dial("tcp", targetAddress)
	if err != nil {
		log.Printf("Failed to connect to target %s: %v\n", targetAddress, err)

		// Respond with a 502 Bad Gateway if connection fails
		response := http.NewResponse(nil)
		response.Protocol = "HTTP/1.1"
		response.Status = 502
		response.Reason = "Bad Gateway"
		response.SetHeader("Content-Type", "text/plain; charset=utf-8")
		response.Body = []byte(fmt.Sprintf("Unable to connect to %s", targetAddress))
		http.WriteResponseToConn(clientConn, response)
		return
	}
	defer targetConn.Close()

	// Respond with 200 Connection Established to the client
	responseLine := "HTTP/1.1 200 Connection Established\r\n\r\n"

	_, err = clientConn.Write([]byte(responseLine))
	if err != nil {
		log.Println("Failed to send connection established response:", err)
		return
	}

	// Create a tunnel: forward data between client and target server
	wg.Add(2)
	forwardTraffic(clientConn, targetConn, &wg)
	wg.Wait()
}

func loadFileLinesAsMap(filename string) map[string]bool {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	lines := make(map[string]bool)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines[line] = true
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatalf("Error reading forbidden hosts file: %v", err)
	}

	return lines
}

func containsBannedWords(body []byte) bool {
	content := string(body)
	for word, _ := range bannedWords {
		if strings.Contains(content, word) {
			return true
		}
	}
	return false
}

// forwardTraffic forwards data between two connections bidirectionally
// To handle bi-directional data forwarding, we'll use io.Copy to copy data from the client to the server and from the server to the client.
func forwardTraffic(clientConn, targetConn net.Conn, wg *sync.WaitGroup) {
	// Create two goroutines to copy data in both directions simultaneously
	go func() {
		_, err := io.Copy(targetConn, clientConn)
		if err != nil {
			log.Println("Error copying data from client to target:", err)
		}
		wg.Done()
	}()

	go func() {
		_, err := io.Copy(clientConn, targetConn)
		if err != nil {
			log.Println("Error copying data from target to client:", err)
		}
		wg.Done()
	}()
}
