package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type ServerConfig struct {
	port string
	host string
}

func main() {
	var sc ServerConfig
	flag.StringVar(&sc.port, "p", "8080", "Port to listen on")
	flag.StringVar(&sc.host, "h", "localhost", "Host to listen on")
	flag.Parse()

	handler := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Received request from", r.RemoteAddr)
		w.Write([]byte("Hello, World!"))
	}
	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", sc.host, sc.port),
		Handler: http.HandlerFunc(handler),
	}
	fmt.Println("Listening on", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
