package server

import (
	"fmt"
	"net"

	"smtp/config"
	"smtp/email"
	"smtp/session"
)

// Server represents the SMTP server
type Server struct {
	config   *config.Config
	listener net.Listener
	store    email.Store
}

// New creates a new SMTP server
func New(cfg *config.Config, store email.Store) *Server {
	return &Server{
		config: cfg,
		store:  store,
	}
}

// Start starts the SMTP server and listens for connections
func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.config.Addr())
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}
	s.listener = ln

	fmt.Printf("📬 SMTP server listening on %s\n", s.config.Addr())

	for {
		conn, err := ln.Accept()
		if err != nil {
			// Check if server is shutting down
			if s.listener == nil {
				return nil
			}
			fmt.Printf("Accept error: %v\n", err)
			continue
		}

		go s.handleConnection(conn)
	}
}

func (s *Server) handleConnection(conn net.Conn) {
	defer conn.Close()

	remoteAddr := conn.RemoteAddr().String()
	fmt.Printf("📥 New connection from %s\n", remoteAddr)

	// Send greeting
	greeting := fmt.Sprintf("220 %s SMTP ready\r\n", s.config.Hostname)
	if _, err := conn.Write([]byte(greeting)); err != nil {
		fmt.Printf("Failed to send greeting: %v\n", err)
		return
	}

	// Create and run session
	sess := session.New(conn, s.store)
	if err := sess.Run(); err != nil {
		fmt.Printf("Session error (%s): %v\n", remoteAddr, err)
	}

	fmt.Printf("📤 Connection closed: %s\n", remoteAddr)
}

// Close shuts down the server
func (s *Server) Close() error {
	if s.listener != nil {
		ln := s.listener
		s.listener = nil
		return ln.Close()
	}
	return nil
}
