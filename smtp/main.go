package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"smtp/config"
	"smtp/email"
	"smtp/server"
)

func main() {
	// Initialize configuration
	cfg := config.DefaultConfig()

	// Initialize email store
	store := email.NewMemoryStore()

	// Create server
	srv := server.New(cfg, store)

	// Handle graceful shutdown
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh

		fmt.Println("\n🛑 Shutting down server...")
		srv.Close()
		os.Exit(0)
	}()

	// Start server
	fmt.Println("🚀 Starting SMTP Server...")
	if err := srv.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
