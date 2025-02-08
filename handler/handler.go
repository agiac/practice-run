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

		msg, err := ParseMessage(string(raw))
		if err != nil {
			member.WriteMessage(fmt.Sprintf("error: bad request: failed to parse message: %v", err))
			continue
		}

		h.handlerMessage(ctx, member, msg)
	}
}

func (h *WebSocketHandler) handlerMessage(ctx context.Context, m *ChatMember, msg interface{}) {
	switch c := msg.(type) {
	case *CreateRoomCommand:
		h.handleCreateRoom(ctx, m, c)
	case *JoinRoomCommand:
		h.handleJoinRoom(ctx, m, c)
	case *LeaveRoomCommand:
		h.handleLeaveRoom(ctx, m, c)
	case *SendMessageCommand:
		h.handleSendMessage(ctx, m, c)
	default:
		log.Printf("Error: unsupported message type: %T", msg)
		m.WriteMessage("server error")
	}
}

func (h *WebSocketHandler) handleCreateRoom(ctx context.Context, m *ChatMember, cmd *CreateRoomCommand) {
	_, err := h.chatService.CreateRoom(ctx, cmd.RoomName)
	if err != nil {
		log.Printf("Debug: %s failed to create room: %v", m.Username(), err)
		m.WriteMessage(fmt.Sprintf("error: failed to create #%s: %v", cmd.RoomName, err))
		return
	}

	m.WriteMessage(fmt.Sprintf("#%s created", cmd.RoomName))
}

func (h *WebSocketHandler) handleJoinRoom(ctx context.Context, m *ChatMember, cmd *JoinRoomCommand) {
	err := h.chatService.AddMember(ctx, cmd.RoomName, m)
	if err != nil {
		log.Printf("Debug: %s failed to join room: %v", m.Username(), err)
		m.WriteMessage(fmt.Sprintf("error: failed to join #%s: %v", cmd.RoomName, err))
		return
	}

	m.WriteMessage(fmt.Sprintf("you've joined #%s", cmd.RoomName))
}

func (h *WebSocketHandler) handleLeaveRoom(ctx context.Context, m *ChatMember, cmd *LeaveRoomCommand) {
	err := h.chatService.RemoveMember(ctx, cmd.RoomName, m)
	if err != nil {
		log.Printf("Debug: %s failed to leave room: %v", m.Username(), err)
		m.WriteMessage(fmt.Sprintf("error: failed to leave #%s: %v", cmd.RoomName, err))
		return
	}

	m.WriteMessage(fmt.Sprintf("you've left #%s", cmd.RoomName))
}

func (h *WebSocketHandler) handleSendMessage(ctx context.Context, m *ChatMember, cmd *SendMessageCommand) {
	err := h.chatService.SendMessage(ctx, cmd.RoomName, m, cmd.Message)
	if err != nil {
		log.Printf("Debug: %s failed to send message: %v", m.Username(), err)
		m.WriteMessage(fmt.Sprintf("error: failed to send message: %v", err))
		return
	}

	m.WriteMessage(fmt.Sprintf("#%s: @%s: %s", cmd.RoomName, m.Username(), cmd.Message))
}
