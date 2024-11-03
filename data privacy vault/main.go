package main

import (
	"context"
	"net/http"
	"vault/handlers"

	"github.com/redis/go-redis/v9"
)

func main() {
	key := []byte("thisis32byteslongkeyforsymmetry!") // 32 bytes for AES-256
	ctx := context.Background()
	redisClient := redis.NewClient(&redis.Options{

		Addr: "localhost:6379",
		// Username: "islamghany",
		Password: "islamghany",
	})

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}

	handler := handlers.NewHandler(key, redisClient)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tokenize", handler.HandleTokenize)
	mux.HandleFunc("POST /detokenize", handler.HandleDetokenize)

	http.ListenAndServe(":8080", mux)
}
