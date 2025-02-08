package main

import (
	"log"
	"net/http"
	"practice-run/provider"
)

func main() {
	if err := http.ListenAndServe(":8080", provider.WebSocketHandler()); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
