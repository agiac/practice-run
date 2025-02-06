package main

import (
	"log"
	"net/http"
	"practice-run/handler"
	"practice-run/service"

	"github.com/gorilla/websocket"
)

func main() {
	h := handler.NewHandler(
		&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		service.NewChat(),
	)

	http.HandleFunc("/ws", h.ServeHTTP)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
