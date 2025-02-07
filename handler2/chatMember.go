package handler2

import (
	"fmt"
	"log"
	"practice-run/room2"
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

func (m *ChatMember) Notify(event room2.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.handleEvent(event)
}

func (m *ChatMember) handleEvent(event room2.Event) {
	eventName := event.Name()

	switch eventName {
	case room2.MessageReceivedEventName:
		m.handleMessageReceivedEvent(event.(*room2.MessageReceivedEvent))
	case room2.MemberJoinedEventName:
		m.handleMemberJoinedEvent(event.(*room2.MemberJoinedEvent))
	case room2.MemberLeftEventName:
		m.handleMemberLeftEvent(event.(*room2.MemberLeftEvent))
	default:
		log.Printf("Error: failed to notify member %s: unknown event %s", m.username, eventName)
	}
}

func (m *ChatMember) handleMessageReceivedEvent(e *room2.MessageReceivedEvent) {
	m.notify(fmt.Sprintf("#%s: @%s: %s", e.RoomName, e.SenderName, e.Message))
}

func (m *ChatMember) handleMemberJoinedEvent(e *room2.MemberJoinedEvent) {
	m.notify(fmt.Sprintf("#%s: @%s joined", e.RoomName, e.MemberName))
}

func (m *ChatMember) handleMemberLeftEvent(e *room2.MemberLeftEvent) {
	m.notify(fmt.Sprintf("#%s: @%s left", e.RoomName, e.MemberName))
}

func (m *ChatMember) notify(message string) {
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to notify member %s: %v", m.username, err)
	}
}
