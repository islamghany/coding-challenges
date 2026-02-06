package main

import (
	"bufio"
	commandsparser "ccmemcached/parser"
	"ccmemcached/store"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
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

	store := store.NewHashTable()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				log.Printf("Temporary error accepting connection %v", netErr)
				continue
			}
			log.Fatalf("Error accepting connection %v", err)
		}
		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store *store.HashTable) {
	defer func() {
		log.Printf("Closing connection from %s\n", conn.RemoteAddr())
		if err := conn.Close(); err != nil {
			log.Printf("Error closing connection %v\n", err)
		}
	}()
	log.Printf("Accepted connection from %s\n", conn.RemoteAddr())

	reader := bufio.NewReader(conn)
	parser := commandsparser.NewParser()

	for {
		fmt.Printf("Waiting for command\n")
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
		fmt.Printf("Command: %+v\n", cmd)
		if cmd.Name == commandsparser.ExitCommand || cmd.Name == commandsparser.EndCommand {
			log.Printf("End command received")
			return
		}

		handleCommand(cmd, store, conn)
	}

}

func handleCommand(cmd *commandsparser.Command, store *store.HashTable, conn net.Conn) {
	switch cmd.Name {
	case commandsparser.SetCommand:
		store.Set(cmd.Key, cmd.Flags, cmd.Expiry, cmd.Value)
		if !cmd.Noreply {
			if _, err := conn.Write([]byte("STORED\r\n")); err != nil {
				log.Printf("Error writing response: %v", err)
			}
		}
	case commandsparser.GetCommand:
		item, ok := store.Get(cmd.Key)
		if !ok {
			if _, err := conn.Write([]byte("NOT_FOUND\r\n")); err != nil {
				log.Printf("Error writing response: %v", err)
			}
			return
		}
		res := fmt.Sprintf("VALUE %s %d %d\r\n", cmd.Key, item.Flags, len(item.Data))
		if _, err := conn.Write([]byte(res)); err != nil {
			log.Printf("Error writing response: %v", err)
		}
		res = string(item.Data) + "\r\n"
		if _, err := conn.Write([]byte(res)); err != nil {
			log.Printf("Error writing response: %v", err)
		}

	case commandsparser.DeleteCommand:
		store.Delete(cmd.Key)
		if !cmd.Noreply {
			if _, err := conn.Write([]byte("DELETED\r\n")); err != nil {
				log.Printf("Error writing response: %v", err)
			}
		}
	default:
		if _, err := conn.Write([]byte("ERROR\r\n")); err != nil {
			log.Printf("Error writing response: %v", err)
		}
	}
}
