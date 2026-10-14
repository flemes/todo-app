package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	addr := getenv("ADDR", ":8080")
	origin := getenv("ALLOWED_ORIGIN", "http://localhost:5173")

	log.Printf("todo-api listening on %s (CORS origin: %s)", addr, origin)
	log.Fatal(http.ListenAndServe(addr, NewRouter(NewStore(), origin)))
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
