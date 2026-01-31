package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

const (
	TCPBufferSize = 4096
	UDPBufferSize = 65535 // Max UDP datagram
)

func main() {
	isUDP := flag.Bool("udp", false, "Use UDP instead of TCP")
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create channel to receive signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		cancel()
	}()

	if *isUDP {
		handleUDPConnection(ctx)
	} else {
		handleTCPConnection(ctx)
	}
}

func handleUDPConnection(ctx context.Context) {
	ln, err := net.ListenPacket("udp", ":7777")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	buf := make([]byte, UDPBufferSize)

	for {
		// Set short deadline to check context
		ln.SetReadDeadline(time.Now().Add(1 * time.Second))

		n, addr, err := ln.ReadFrom(buf)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				select {
				case <-ctx.Done():
					log.Println("UDP server shutting down")
					return
				default:
					continue
				}
			}
			log.Printf("ReadFrom error: %v", err)
			continue
		}

		ln.WriteTo(buf[:n], addr)
	}
}

func handleTCPConnection(ctx context.Context) {
	var wg sync.WaitGroup

	ln, err := net.Listen("tcp", ":7777")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	// Goroutine to close listener when context cancelled
	go func() {
		<-ctx.Done()
		log.Println("Closing listener...")
		ln.Close()
	}()

	for {

		conn, err := ln.Accept()
		if err != nil {
			// Check if it's because we're shutting down
			select {
			case <-ctx.Done():
				log.Println("Stopped accepting connections")
				wg.Wait() // ← Wait for active connections!
				log.Println("All connections closed")
				return
			default:
				log.Printf("Accept error: %v", err)
				continue
			}
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			handleConnection(ctx, conn)
		}()
	}
}

func handleConnection(ctx context.Context, conn net.Conn) {
	defer conn.Close()

	log.Println("Accepted connection from ", conn.RemoteAddr())
	buf := make([]byte, TCPBufferSize)
	for {
		conn.SetReadDeadline(time.Now().Add(20 * time.Second))
		// conn.Read will block until data is available or the read deadline is reached
		n, err := conn.Read(buf)
		if err != nil {

			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// Timeout - check if we should shutdown
				select {
				case <-ctx.Done():
					log.Println("Connection closing due to shutdown")
					return
				default:
					continue // Keep waiting for data
				}
			}
			// Real error or EOF
			if err != io.EOF {
				log.Printf("Read error: %v", err)
			}
			return
		}

		log.Println("Received message: ", buf[:n])
		_, err = conn.Write(buf[:n])
		if err != nil {
			log.Printf("Error writing message: %v", err)
			return
		}

	}
}
