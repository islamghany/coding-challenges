package main

import (
	"fmt"
	"webserver/http"
	"webserver/http/handler"
	"webserver/http/request"
	"webserver/http/response"
	"webserver/http/router"
)

func main() {

	router := router.NewRouter()
	router.Get("/v1/hello", handler.HandlerFunc(func(w response.Response, r *request.Request) {
		fmt.Print("Hello World")
		w.Write([]byte("Hello World"))
	}))
	router.Post("/v1/hello", handler.HandlerFunc(func(w response.Response, r *request.Request) {
		fmt.Print("Hello World")
		w.Write([]byte("Hello World POST"))
	}))
	router.Put("/v1/hello/:id", handler.HandlerFunc(func(w response.Response, r *request.Request) {
		fmt.Print("Hello World")
		w.Write([]byte("Hello World PUT"))
	}))
	app := http.NewServer(http.ServerConfig{
		Addr:   ":8080",
		Router: router,
	})

	fmt.Println("Server is running on port 8080")
	app.Start()
}
