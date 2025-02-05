package handler

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Handler struct {
	u websocket.Upgrader
}

func NewHandler() *Handler {
	return &Handler{
		u: websocket.Upgrader{}, // TODO: configure
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := h.u.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error: failed to upgrade connection: %v", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			log.Printf("Error: failed to close connection: %v", err)
		}
	}(conn)

	err = conn.WriteMessage(websocket.TextMessage, []byte("create a room"))
	if err != nil {
		log.Printf("Error: failed to write message: %v", err)
		return
	}

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error: failed to read message: %v", err)
			break
		}
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Printf("Error: failed to write message: %v", err)
			break
		}
	}
}
