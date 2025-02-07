package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"practice-run/room"

	"github.com/gorilla/websocket"
)

//go:generate mockgen -destination mocks/chat_service_mock.go -package mocks . ChatService
type ChatService interface {
	AddMemberToRoom(ctx context.Context, roomName string, member room.Member) error
	RemoveMemberFromRoom(ctx context.Context, roomName string, member room.Member) error
	SendMessageToRoom(ctx context.Context, roomName string, member room.Member, message string) error
}

type WebSocketHandler struct {
	upgrader    *websocket.Upgrader
	chatService ChatService
}

func NewWebSocketHandler(upgrader *websocket.Upgrader, chatService ChatService) *WebSocketHandler {
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
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error: failed to read message, breaking connection: %v", err)
			break
		}

		if mt != websocket.TextMessage {
			member.WriteMessage("bad request: only text messages are supported")
			continue
		}

		cmd, err := ParseMessage(string(msg))
		if err != nil {
			member.WriteMessage(fmt.Sprintf("bad request: failed to parse message: %v", err))
			continue
		}

		h.handleCommand(ctx, member, cmd)
	}
}

func (h *WebSocketHandler) handleCommand(ctx context.Context, m *WSChatMember, cmd interface{}) {
	switch c := cmd.(type) {
	case *JoinChannelCommand:
		h.handleJoinChannel(ctx, m, c)
	case *LeaveChannelCommand:
		h.handleLeaveChannel(ctx, m, c)
	case *SendMessageCommand:
		h.handleSendMessage(ctx, m, c)
	default:
		log.Printf("Error: unsupported command type: %T", cmd)
		m.WriteMessage("server error")
	}
}

func (h *WebSocketHandler) handleJoinChannel(ctx context.Context, m *WSChatMember, cmd *JoinChannelCommand) {
	err := h.chatService.AddMemberToRoom(ctx, cmd.ChannelName, m)
	if err != nil {
		log.Printf("Debug: %s failed to join channel: %v", m.Username(), err)
		m.WriteMessage(fmt.Sprintf("failed to join #%s: %v", cmd.ChannelName, err))
		return
	}

	m.WriteMessage(fmt.Sprintf("you've joined #%s", cmd.ChannelName))
}

func (h *WebSocketHandler) handleLeaveChannel(ctx context.Context, m *WSChatMember, cmd *LeaveChannelCommand) {
	err := h.chatService.RemoveMemberFromRoom(ctx, cmd.ChannelName, m)
	if err != nil {
		log.Printf("Debug: %s failed to leave channel: %v", m.Username(), err)
		m.WriteMessage(fmt.Sprintf("failed to leave #%s: %v", cmd.ChannelName, err))
		return
	}

	m.WriteMessage(fmt.Sprintf("you've left #%s", cmd.ChannelName))
}

func (h *WebSocketHandler) handleSendMessage(ctx context.Context, m *WSChatMember, cmd *SendMessageCommand) {
	err := h.chatService.SendMessageToRoom(ctx, cmd.ChannelName, m, cmd.Message)
	if err != nil {
		log.Printf("Debug: %s failed to send message: %v", m.Username(), err)
		m.WriteMessage(fmt.Sprintf("failed to send message: %v", err))
		return
	}

	m.WriteMessage(fmt.Sprintf("#%s: @%s: %s", cmd.ChannelName, m.Username(), cmd.Message))
}
