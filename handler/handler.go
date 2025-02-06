package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

type WSHandler interface {
	ServeWS(ctx context.Context, msg io.Reader, conn *websocket.Conn)
}

type WSServer struct {
	u websocket.Upgrader
	h WSHandler
}

func NewWSServer(u websocket.Upgrader, h WSHandler) *WSServer {
	return &WSServer{
		u: u,
		h: h,
	}
}

var userIdCounter = atomic.Int32{}

func (h *WSServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Assign identifier to the user
	userId := userIdCounter.Add(1)

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
				continue
			}

			err = h.handleCreateRoom(crm)
			if err != nil {
				_ = NewErrorMessage(err.Error()).Send(conn)
				continue
			}

			_ = NewInfoMessage(fmt.Sprintf("room %s created", crm.RoomName)).Send(conn)
		case JoinRoom:
			var jrm JoinRoomMessageBody
			err = json.Unmarshal(msg.Body, &jrm)
			if err != nil {
				_ = NewErrorMessage("failed to unmarshal message").Send(conn)
				continue
			}

			err = h.handleJoinRoom(jrm)
			if err != nil {
				_ = NewErrorMessage(err.Error()).Send(conn)
				continue
			}

			_ = NewInfoMessage(fmt.Sprintf("room %s joined", jrm.RoomName)).Send(conn)
		case LeaveRoom:
			var lrm LeaveRoomMessageBody
			err = json.Unmarshal(msg.Body, &lrm)
			if err != nil {
				_ = NewErrorMessage("failed to unmarshal message").Send(conn)
				continue
			}

			err = h.handleLeaveRoom(lrm)
			if err != nil {
				_ = NewErrorMessage(err.Error()).Send(conn)
				continue
			}

			_ = NewInfoMessage(fmt.Sprintf("room %s left", lrm.RoomName)).Send(conn)
		case SendMessageToRoom:
			var smrm SendMessageMessageBody
			err = json.Unmarshal(msg.Body, &smrm)
			if err != nil {
				_ = NewErrorMessage("failed to unmarshal message").Send(conn)
				continue
			}

			err = h.handleSendMessageToRoom(smrm)
			if err != nil {
				_ = NewErrorMessage(err.Error()).Send(conn)
				continue
			}

			_ = NewInfoMessage(fmt.Sprintf("message sent to room %s", smrm.RoomName)).Send(conn)
		}
	}
}

func (h *WSServer) handleCreateRoom(msg CreateRoomMessageBody) error {
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

func (h *WSServer) handleJoinRoom(msg JoinRoomMessageBody) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.db.Rooms[msg.RoomName]
	if !ok {
		return fmt.Errorf("room does not exist")
	}

	_, ok = room.Participants[msg.ParticipantName]
	if ok {
		return fmt.Errorf("participant already in the room")
	}

	room.Participants[msg.ParticipantName] = Participant{
		Name: msg.ParticipantName,
		Conn: nil,
	}

	return nil
}

func (h *WSServer) handleLeaveRoom(lrm LeaveRoomMessageBody) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	room, ok := h.db.Rooms[lrm.RoomName]
	if !ok {
		return fmt.Errorf("room does not exist")
	}

	_, ok = room.Participants[lrm.ParticipantName]
	if !ok {
		return fmt.Errorf("participant not in the room")
	}

	delete(room.Participants, lrm.ParticipantName)

	return nil
}
