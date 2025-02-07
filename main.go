package main

import (
	"log"
	"net/http"
	"practice-run/handler"
	"practice-run/service"

	"github.com/gorilla/websocket"
)

func main() {
	// TODO: improve context and shutdown handling

	u := &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	s := service.NewChat()

	h := handler.NewWebSocketHandler(u, s)

	http.HandleFunc("/ws", h.ServeHTTP)

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
