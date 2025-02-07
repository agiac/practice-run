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
	u *websocket.Upgrader
	s ChatService
}

func NewWebSocketHandler(u *websocket.Upgrader, s ChatService) *WebSocketHandler {
	return &WebSocketHandler{u: u, s: s}
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

	member := NewChatMember(username, conn)

	member.Notify(fmt.Sprintf("welcome, %s!", username))

	for {
		mt, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error: failed to read message, breaking connection: %v", err)
			break
		}

		if mt != websocket.TextMessage {
			member.Notify("bad request: only text messages are supported")
			continue
		}

		msg, err := ParseMessage(string(raw))
		if err != nil {
			member.Notify(fmt.Sprintf("bad request: failed to parse message: %v", err))
			continue
		}

		switch cmd := msg.(type) {
		case *JoinChannelCommand:
			h.handleJoinChannel(ctx, member, cmd)
		case *LeaveChannelCommand:
			h.handleLeaveChannel(ctx, member, cmd)
		case *SendMessageCommand:
			h.handleSendMessage(ctx, member, cmd)
		//case *SendDirectMessageCommand:
		//case *ListChannelsCommand:
		//case *ListChannelUsersCommand:
		default:
			log.Printf("Error: unsupported message type: %T", cmd)
			member.Notify("server error")
		}
	}
}

func (h *WebSocketHandler) handleJoinChannel(ctx context.Context, m *WSChatMember, cmd *JoinChannelCommand) {
	err := h.s.AddMemberToRoom(ctx, cmd.ChannelName, m)
	if err != nil {
		log.Printf("Debug: %s failed to join channel: %v", m.Username(), err)
		m.Notify(fmt.Sprintf("failed to join channel: %v", err))
		return
	}
	m.Notify(fmt.Sprintf("%s joined channel #%s", m.Username(), cmd.ChannelName))
}

func (h *WebSocketHandler) handleLeaveChannel(ctx context.Context, m *WSChatMember, cmd *LeaveChannelCommand) {
	err := h.s.RemoveMemberFromRoom(ctx, cmd.ChannelName, m)
	if err != nil {
		log.Printf("Debug: %s failed to leave channel: %v", m.Username(), err)
		m.Notify(fmt.Sprintf("failed to leave channel: %v", err))
		return
	}

	m.Notify(fmt.Sprintf("%s left channel #%s", m.Username(), cmd.ChannelName))
}

func (h *WebSocketHandler) handleSendMessage(ctx context.Context, m *WSChatMember, cmd *SendMessageCommand) {
	err := h.s.SendMessageToRoom(ctx, cmd.ChannelName, m, cmd.Message)
	if err != nil {
		log.Printf("Debug: %s failed to send message: %v", m.Username(), err)
		m.Notify(fmt.Sprintf("failed to send message: %v", err))
		return
	}
	m.Notify(fmt.Sprintf("#%s: @%s: %s", cmd.ChannelName, m.Username(), cmd.Message))
}
