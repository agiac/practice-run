package main

import (
	"log"
	"net/http"
	"practice-run/handler2"
	"practice-run/service"

	"github.com/gorilla/websocket"
)

func main() {

	h := handler2.NewHandler(
		&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		service.NewChat(),
	)

	http.HandleFunc("/ws", h.ServeHTTP)

	http.HandleFunc("/", http.FileServer(http.Dir("./public")).ServeHTTP)

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
