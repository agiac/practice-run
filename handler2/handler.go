package handler2

import (
	"fmt"
	"log"
	"net/http"
	"practice-run/messages"

	"github.com/gorilla/websocket"
)

type Handler struct {
	u *websocket.Upgrader
}

func NewHandler(u *websocket.Upgrader) *Handler {
	return &Handler{u: u}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read credentials, skip validation
	username, _, ok := r.BasicAuth()
	if !ok {
		log.Printf("Debug: unauthorized request")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	conn, err := h.u.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error: failed to upgrade connection: %v", err)
		return
	}

	defer func(conn *websocket.Conn) {
		err = conn.Close()
		if err != nil {
			log.Printf("Error: failed to close connection: %v", err)
		}
	}(conn)

	// Send welcome message upon successful connection
	h.WriteMessage(conn, fmt.Sprintf("Welcome, %s!", username))

	// Handle messages
	for {
		mt, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error: failed to read message, breaking connection: %v", err)
			break
		}

		if mt != websocket.TextMessage {
			h.WriteMessage(conn, "bad request: only text messages are supported")
			continue
		}

		// Handle message
		msg, err := messages.ParseMessage(string(raw))
		if err != nil {
			h.WriteMessage(conn, fmt.Sprintf("bad request: failed to parse message: %v", err))
			continue
		}

		// Type switch
		switch m := msg.(type) {
		case *messages.JoinChannelCommand:
			h.WriteMessage(conn, fmt.Sprintf("%s joined channel #%s", username, m.ChannelName))
		// TODO: implement
		case *messages.LeaveChannelCommand:
			// TODO: implement
		case *messages.SendMessageCommand:
			// TODO: implement
		case *messages.SendDirectMessageCommand:
		// TODO: implement
		case *messages.ListChannelsCommand:
			// TODO: implement
		case *messages.ListChannelUsersCommand:
			// TODO: implement
		default:
			log.Printf("Error: unsupported message type: %T", m)
			h.WriteMessage(conn, "server error")
		}
	}
}

func (h *Handler) WriteMessage(conn *websocket.Conn, msg string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("Error: failed to write message: %v", err)
	}
}
