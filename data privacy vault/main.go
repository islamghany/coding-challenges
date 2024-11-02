package main

import (
	"net/http"
	"valut/handlers"
)

func main() {
	key := []byte("thisis32byteslongkeyforsymmetry!") // 32 bytes for AES-256
	// fmt.Println(cryptographic.SHA256Hash("Hello World"))
	handler := handlers.NewHandler(key)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /tokenize", handler.HandleTokenize)
	mux.HandleFunc("POST /detokenize", handler.HandleDetokenize)

	http.ListenAndServe(":8080", mux)
}
