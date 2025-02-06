package handler2

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

//go:generate mockgen -destination mocks/chat_service_mock.go -package mocks . ChatService
type ChatService interface {
	JoinChannel(ctx context.Context, username, channelName string) error
	LeaveChannel(ctx context.Context, username, channelName string) error
	SendMessage(ctx context.Context, username, channelName, message string) error
	GetUpdates(ctx context.Context, username string) (<-chan string, error)
}

type Handler struct {
	u *websocket.Upgrader
	s ChatService
}

func NewHandler(u *websocket.Upgrader, s ChatService) *Handler {
	return &Handler{u: u, s: s}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	// Send welcome message upon successful connection
	h.WriteMessage(conn, fmt.Sprintf("welcome, %s!", username))

	// Get updates
	updates, err := h.s.GetUpdates(ctx, username)
	if err != nil {
		log.Printf("Error: failed to get updates stream for user %s: %v", username, err)
		return
	}

	go func() {
		for {
			select {
			case update := <-updates:
				h.WriteMessage(conn, update)
			case <-ctx.Done():
				return
			}
		}
	}()

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
		msg, err := ParseMessage(string(raw))
		if err != nil {
			h.WriteMessage(conn, fmt.Sprintf("bad request: failed to parse message: %v", err))
			continue
		}

		switch m := msg.(type) {
		case *JoinChannelCommand:
			h.handleJoinChannel(ctx, conn, username, m)
		case *LeaveChannelCommand:
			h.handleLeaveChannel(ctx, conn, username, m)
		case *SendMessageCommand:
			h.handleSendMessage(ctx, conn, username, m)
		//case *SendDirectMessageCommand:
		//case *ListChannelsCommand:
		//case *ListChannelUsersCommand:
		default:
			log.Printf("Error: unsupported message type: %T", m)
			h.WriteMessage(conn, "server error")
		}
	}
}

func (h *Handler) handleJoinChannel(ctx context.Context, conn *websocket.Conn, username string, cmd *JoinChannelCommand) {
	err := h.s.JoinChannel(ctx, username, cmd.ChannelName)
	if err != nil {
		log.Printf("Debug: %s failed to join channel: %v", username, err)
		h.WriteMessage(conn, fmt.Sprintf("failed to join channel: %v", err))
		return
	}
	h.WriteMessage(conn, fmt.Sprintf("%s joined channel #%s", username, cmd.ChannelName))
}

func (h *Handler) handleLeaveChannel(ctx context.Context, conn *websocket.Conn, username string, cmd *LeaveChannelCommand) {
	err := h.s.LeaveChannel(ctx, username, cmd.ChannelName)
	if err != nil {
		log.Printf("Debug: %s failed to leave channel: %v", username, err)
		h.WriteMessage(conn, fmt.Sprintf("failed to leave channel: %v", err))
		return
	}
	h.WriteMessage(conn, fmt.Sprintf("%s left channel #%s", username, cmd.ChannelName))
}

func (h *Handler) handleSendMessage(ctx context.Context, conn *websocket.Conn, username string, cmd *SendMessageCommand) {
	err := h.s.SendMessage(ctx, username, cmd.ChannelName, cmd.Message)
	if err != nil {
		log.Printf("Debug: %s failed to send message: %v", username, err)
		h.WriteMessage(conn, fmt.Sprintf("failed to send message: %v", err))
		return
	}
	h.WriteMessage(conn, fmt.Sprintf("#%s: @%s: %s", cmd.ChannelName, username, cmd.Message))
}

func (h *Handler) WriteMessage(conn *websocket.Conn, msg string) {
	err := conn.WriteMessage(websocket.TextMessage, []byte(msg))
	if err != nil {
		log.Printf("Error: failed to write message: %v", err)
	}
}
