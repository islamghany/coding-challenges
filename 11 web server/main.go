package main

import (
	"fmt"
	"webserver/http"
)

func main() {
	app := http.NewServer(http.ServerConfig{
		Addr: ":8080",
	})

	app.Get("/v1/hello", func(w http.Response, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	fmt.Println("Server is running on port 8080")
	app.Start()
}
