package main

import (
	// "database/sql"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	_ "github.com/mattn/go-sqlite3"
	"github.com/spaolacci/murmur3"
)

// Define an envelope type.
type envelope map[string]interface{}

func main() {
	// // Initialize the database
	fmt.Println("Initializing database")
	db, err := initDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	mux := http.NewServeMux()

	mux.HandleFunc("POST /shorten", func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			URL string `json:"url"`
		}
		var req response
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// Generate a hash for the URL
		hash := Hash(req.URL)
		// Encode the hash using base62
		shortened := EncodeBase62(hash)
		// Insert the URL and the shortened version into the database
		_, err = db.Exec("INSERT INTO urls_shortened (url, shortened) VALUES (?, ?)", req.URL, shortened)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// fmt.Fprintf(w, "Shortened URL: %s\n", shortened)

		// Return the shortened URL
		w.Header().Set("Content-Type", "application/json")
		shortened = fmt.Sprintf("http://localhost:8080/%s", shortened)
		data := map[string]string{"shortened": shortened}
		jsonData, err := json.MarshalIndent(envelope{"data": data}, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonData = append(jsonData, '\n')
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	mux.HandleFunc("GET /{shortened}", func(w http.ResponseWriter, r *http.Request) {
		// Get the shortened URL from the request
		shortened := r.URL.Path[1:]
		// Query the database for the original URL
		var original string
		err := db.QueryRow("SELECT url FROM urls_shortened WHERE shortened = ?", shortened).Scan(&original)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Redirect the user to the original URL
		http.Redirect(w, r, original, http.StatusMovedPermanently)
	})

	// DELETE /shorten/{shortened}
	mux.HandleFunc("DELETE /shorten/{shortened}", func(w http.ResponseWriter, r *http.Request) {
		// Get the shortened URL from the request
		shortened := r.URL.Path[1:]
		// Query the database for the original URL
		var original string
		err := db.QueryRow("SELECT url FROM urls_shortened WHERE shortened = ?", shortened).Scan(&original)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		// Delete the URL from the database
		_, err = db.Exec("DELETE FROM urls_shortened WHERE shortened = ?", shortened)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Return a success message
		w.Header().Set("Content-Type", "application/json")
		data := map[string]string{"message": "URL deleted successfully"}
		jsonData, err := json.MarshalIndent(envelope{"data": data}, "", "\t")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonData = append(jsonData, '\n')
		w.WriteHeader(http.StatusOK)
		w.Write(jsonData)
	})

	srv := &http.Server{
		Addr:    ":8080",
		Handler: CorsMiddleware(mux),
	}
	fmt.Println("Server is running on port 8080")
	if err := srv.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func initDB() (*sql.DB, error) {
	// Open a connection to the SQLite database
	// If the file does not exist, it will be created.
	fmt.Println("Opening database")
	database, err := sql.Open("sqlite3", "./urls.db")
	if err != nil {
		return nil, err
	}

	// Create the table if it does not exist
	urlsShortened := `
	CREATE TABLE IF NOT EXISTS urls_shortened (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		url TEXT NOT NULL,
		shortened TEXT NOT NULL
	);
	`
	fmt.Println("Creating table")
	_, err = database.Exec(urlsShortened)
	if err != nil {
		return nil, err
	}
	fmt.Println("Table created successfully")
	return database, nil
}

func Hash(input string) uint64 {

	// Generate a 64-bit hash using MurmurHash3
	hash64 := murmur3.Sum64([]byte(input))

	return hash64
}

var base62Table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func EncodeBase62(num uint64) string {

	if num == 0 {
		return string(base62Table[0])
	}
	var encoded strings.Builder
	for num > 0 {
		rem := num % 62
		encoded.WriteByte(base62Table[rem])
		num = num / 62
	}

	return reverse(encoded.String())
}

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			return
		}
		next.ServeHTTP(w, r)
	})
}
