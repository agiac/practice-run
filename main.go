package main

import (
	"log"
	"net/http"
	"practice-run/provider2"
)

func main() {
	if err := http.ListenAndServe(":8080", provider2.WebSocketHandler()); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
