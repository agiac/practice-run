package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"practice-run/chat"

	"github.com/gorilla/websocket"
)

//go:generate mockgen -destination mocks/chat_service_mock.go -mock_names chatService=ChatService -package mocks . chatService
type chatService interface {
	CreateRoom(ctx context.Context, roomName string) (*chat.Room, error)
	AddMember(ctx context.Context, roomName string, member chat.Member) error
	RemoveMember(ctx context.Context, roomName string, member chat.Member) error
	SendMessage(ctx context.Context, roomName string, member chat.Member, message string) error
}

type WebSocketHandler struct {
	upgrader    *websocket.Upgrader
	chatService chatService
}

func NewWebSocketHandler(upgrader *websocket.Upgrader, chatService chatService) *WebSocketHandler {
	return &WebSocketHandler{upgrader: upgrader, chatService: chatService}
}

func (h *WebSocketHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	username, _, ok := r.BasicAuth()
	if !ok {
		log.Printf("Debug: unauthorized request")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	// skip authentication for now

	conn, err := h.upgrader.Upgrade(w, r, nil)
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

	member := NewChatMember(username, conn)

	for {
		mt, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error: failed to read message, breaking connection: %v", err)
			break
		}

		if mt != websocket.TextMessage {
			member.WriteMessage("error: bad request: only text messages are supported")
			continue
		}

		cmd, err := ParseMessage(string(raw))
		if err != nil {
			member.WriteMessage(fmt.Sprintf("error: bad request: failed to parse message: %v", err))
			continue
		}

		err = cmd.Execute(ctx, member, h.chatService)
		if err != nil {
			log.Printf("Debug: failed to execute command %s %+v: %v", cmd.Name(), cmd, err)
			member.WriteMessage(fmt.Sprintf("error: %v", err))
		}
	}
}
