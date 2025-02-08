package handler

import (
	"fmt"
	"log"
	"practice-run/chat"
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
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to write message to member %s: %v", m.username, err)
	}
}

func (m *ChatMember) Notify(event chat.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.handleEvent(event)
}

func (m *ChatMember) handleEvent(event chat.Event) {
	eventName := event.Name()

	switch eventName {
	case chat.MessageReceivedEventName:
		m.handleMessageReceivedEvent(event.(*chat.MessageReceivedEvent))
	case chat.MemberJoinedEventName:
		m.handleMemberJoinedEvent(event.(*chat.MemberJoinedEvent))
	case chat.MemberLeftEventName:
		m.handleMemberLeftEvent(event.(*chat.MemberLeftEvent))
	default:
		log.Printf("Error: failed to notify member %s: unknown event %s", m.username, eventName)
	}
}

func (m *ChatMember) handleMessageReceivedEvent(e *chat.MessageReceivedEvent) {
	m.notify(fmt.Sprintf("#%s: @%s: %s", e.RoomName, e.SenderName, e.Message))
}

func (m *ChatMember) handleMemberJoinedEvent(e *chat.MemberJoinedEvent) {
	m.notify(fmt.Sprintf("#%s: @%s joined", e.RoomName, e.MemberName))
}

func (m *ChatMember) handleMemberLeftEvent(e *chat.MemberLeftEvent) {
	m.notify(fmt.Sprintf("#%s: @%s left", e.RoomName, e.MemberName))
}

func (m *ChatMember) notify(message string) {
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to notify member %s: %v", m.username, err)
	}
}
