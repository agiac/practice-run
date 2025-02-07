package main

import (
	"log"
	"net/http"
	"os"
	"practice-run/internal/provider"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}

	if err := http.ListenAndServe(port, provider.WebSocketHandler()); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
