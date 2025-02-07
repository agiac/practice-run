package main

import (
	"log"
	"net/http"
	"practice-run/handler"
)

func main() {
	// TODO: improve context and shutdown handling

	err := http.ListenAndServe(":8080", handler.MakeHandler())

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
