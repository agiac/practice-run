package main

import (
	"fmt"
	"log"
	"net/http"
	"practice-run/provider"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, "Practice Run")
	})

	http.Handle("/ws", provider.WebSocketHandler())

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
