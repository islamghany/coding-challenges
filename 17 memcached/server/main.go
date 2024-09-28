package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type ServerConfig struct {
	Port int
}

func parseConfig() ServerConfig {
	var config ServerConfig
	flag.IntVar(&config.Port, "p", 11211, "Port to listen on")
	flag.Parse()
	return config
}

func main() {
	config := parseConfig()
	log.Printf("Starting server on port %d\n", config.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", config.Port))
	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}
	defer listener.Close()

	// Setup Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-stop
		log.Println("Shutting down server...")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
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
	defer conn.Close()
	conn.Write([]byte("Hello, world!\n"))
}
