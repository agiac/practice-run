package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"practice-run/chat"

	"github.com/gorilla/websocket"
)

//go:generate mockgen -destination mocks/chat_service_mock.go -package mocks . ChatService
type ChatService interface {
	JoinChannel(ctx context.Context, channelName string, member chat.Member) error
	LeaveChannel(ctx context.Context, channelName string, member chat.Member) error
	SendMessage(ctx context.Context, channelName string, member chat.Member, message string) error
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

	// Create channel member
	member := chat.NewChatMember(username, conn)

	// Send welcome message upon successful connection
	member.Notify(fmt.Sprintf("welcome, %s!", username))
	//h.WriteMessage(conn, fmt.Sprintf("welcome, %s!", username))

	// Get updates
	//updates, err := h.s.GetUserUpdates(ctx, username)
	//if err != nil {
	//	log.Printf("Error: failed to get updates stream for user %s: %v", username, err)
	//	return
	//}
	//
	//go func() {
	//	for {
	//		select {
	//		case update := <-updates:
	//			h.WriteMessage(conn, update)
	//		case <-ctx.Done():
	//			return
	//		}
	//	}
	//}()

	// Handle messages
	for {
		mt, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("Error: failed to read message, breaking connection: %v", err)
			break
		}

		if mt != websocket.TextMessage {
			//h.WriteMessage(conn, "bad request: only text messages are supported")
			member.Notify("bad request: only text messages are supported")
			continue
		}

		// Handle message
		msg, err := ParseMessage(string(raw))
		if err != nil {
			//h.WriteMessage(conn, fmt.Sprintf("bad request: failed to parse message: %v", err))
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
			//h.WriteMessage(conn, "server error")
			member.Notify("server error")
		}
	}
}

func (h *WebSocketHandler) handleJoinChannel(ctx context.Context, m *chat.WSChatMember, cmd *JoinChannelCommand) {
	err := h.s.JoinChannel(ctx, cmd.ChannelName, m)
	if err != nil {
		log.Printf("Debug: %s failed to join channel: %v", m.Username(), err)
		//h.WriteMessage(conn, fmt.Sprintf("failed to join channel: %v", err))
		m.Notify(fmt.Sprintf("failed to join channel: %v", err))
		return
	}
	//h.WriteMessage(conn, fmt.Sprintf("%s joined channel #%s", username, cmd.ChannelName))
	m.Notify(fmt.Sprintf("%s joined channel #%s", m.Username(), cmd.ChannelName))
}

func (h *WebSocketHandler) handleLeaveChannel(ctx context.Context, m *chat.WSChatMember, cmd *LeaveChannelCommand) {
	err := h.s.LeaveChannel(ctx, cmd.ChannelName, m)
	if err != nil {
		log.Printf("Debug: %s failed to leave channel: %v", m.Username(), err)
		//h.WriteMessage(conn, fmt.Sprintf("failed to leave channel: %v", err))
		m.Notify(fmt.Sprintf("failed to leave channel: %v", err))
		return
	}
	//h.WriteMessage(conn, fmt.Sprintf("%s left channel #%s", username, cmd.ChannelName))
	m.Notify(fmt.Sprintf("%s left channel #%s", m.Username(), cmd.ChannelName))
}

func (h *WebSocketHandler) handleSendMessage(ctx context.Context, m *chat.WSChatMember, cmd *SendMessageCommand) {
	err := h.s.SendMessage(ctx, cmd.ChannelName, m, cmd.Message)
	if err != nil {
		log.Printf("Debug: %s failed to send message: %v", m.Username(), err)
		//h.WriteMessage(conn, fmt.Sprintf("failed to send message: %v", err))
		m.Notify(fmt.Sprintf("failed to send message: %v", err))
		return
	}
	//h.WriteMessage(conn, fmt.Sprintf("#%s: @%s: %s", cmd.ChannelName, username, cmd.Message))
	m.Notify(fmt.Sprintf("#%s: @%s: %s", cmd.ChannelName, m.Username(), cmd.Message))
}
