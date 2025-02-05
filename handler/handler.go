package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Participant struct {
	Name string
	Conn websocket.Conn
}

type Room struct {
	Name         string
	Participants map[string]Participant
}

type DB struct {
	Rooms map[string]Room
}

type Handler struct {
	u  websocket.Upgrader
	mu sync.Mutex
	db *DB
}

func NewHandler() *Handler {
	return &Handler{
		u: websocket.Upgrader{}, // TODO: configure
		db: &DB{
			Rooms: make(map[string]Room),
		},
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

	for {
		mt, reader, err := conn.NextReader()
		if err != nil {
			log.Printf("Error: failed to read message: %v", err)
			break
		}

		if mt != websocket.TextMessage {
			_ = NewErrorMessage("only text messages are supported").Send(conn)
			break
		}

		var msg GenericMessage
		err = json.NewDecoder(reader).Decode(&msg)
		if err != nil {
			_ = NewErrorMessage("failed to unmarshal message").Send(conn)
		}

		switch msg.Type {
		case CreateRoom:
			var crm CreateRoomMessageBody
			err = json.Unmarshal(msg.Body, &crm)
			if err != nil {
				_ = NewErrorMessage("failed to unmarshal message").Send(conn)
			}

			err = h.handleCreateRoom(crm)
			if err != nil {
				_ = NewErrorMessage(err.Error()).Send(conn)
			}

			_ = NewInfoMessage(fmt.Sprintf("room %s created", crm.RoomName)).Send(conn)
		}
	}
}

func (h *Handler) handleCreateRoom(msg CreateRoomMessageBody) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.db.Rooms[msg.RoomName]; ok {
		return fmt.Errorf("room already exists")
	}

	h.db.Rooms[msg.RoomName] = Room{
		Name:         msg.RoomName,
		Participants: make(map[string]Participant),
	}

	return nil
}
