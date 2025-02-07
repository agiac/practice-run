package main

import (
	"log"
	"net/http"
	"practice-run/internal/provider"
)

func main() {
	if err := http.ListenAndServe(":8080", provider.WebSocketHandler()); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
