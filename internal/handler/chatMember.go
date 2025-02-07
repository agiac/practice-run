package handler

import (
	"fmt"
	"log"
	"practice-run/internal/room"
	"sync"

	"github.com/gorilla/websocket"
)

type ChatMember struct {
	mu   sync.Mutex
	conn *websocket.Conn

	username string
}

func NewChatMember(username string, conn *websocket.Conn) *ChatMember {
	return &ChatMember{username: username, mu: sync.Mutex{}, conn: conn}
}

func (m *ChatMember) Username() string {
	return m.username
}

func (m *ChatMember) WriteMessage(message string) {
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to write message to member %s: %v", m.username, err)
	}
}

func (m *ChatMember) Notify(event room.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.handleEvent(event)
}

func (m *ChatMember) handleEvent(event room.Event) {
	eventName := event.Name()

	switch eventName {
	case room.MessageReceivedEventName:
		m.handleMessageReceivedEvent(event.(*room.MessageReceivedEvent))
	case room.MemberJoinedEventName:
		m.handleMemberJoinedEvent(event.(*room.MemberJoinedEvent))
	case room.MemberLeftEventName:
		m.handleMemberLeftEvent(event.(*room.MemberLeftEvent))
	default:
		log.Printf("Error: failed to notify member %s: unknown event %s", m.username, eventName)
	}
}

func (m *ChatMember) handleMessageReceivedEvent(e *room.MessageReceivedEvent) {
	m.notify(fmt.Sprintf("#%s: @%s: %s", e.RoomName, e.SenderName, e.Message))
}

func (m *ChatMember) handleMemberJoinedEvent(e *room.MemberJoinedEvent) {
	m.notify(fmt.Sprintf("#%s: @%s joined", e.RoomName, e.MemberName))
}

func (m *ChatMember) handleMemberLeftEvent(e *room.MemberLeftEvent) {
	m.notify(fmt.Sprintf("#%s: @%s left", e.RoomName, e.MemberName))
}

func (m *ChatMember) notify(message string) {
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to notify member %s: %v", m.username, err)
	}
}
