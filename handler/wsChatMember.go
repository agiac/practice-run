package handler

import (
	"fmt"
	"log"
	"practice-run/room"
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

func (m *ChatMember) Notify(event room.Event) {
	m.mu.Lock()
	defer m.mu.Unlock()

	eventName := event.Name()

	switch eventName {
	case room.MessageReceivedEventName:
		e := event.(*room.MessageReceivedEvent)
		m.notify(fmt.Sprintf("#%s: @%s: %s", e.RoomName, e.SenderName, e.Message))
	default:
		log.Printf("Error: failed to notify member %s: unknown event %s", m.username, eventName)
	}
}

func (m *ChatMember) notify(message string) {
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to notify member %s: %v", m.username, err)
	}
}

func (m *ChatMember) WriteMessage(message string) {
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to write message to member %s: %v", m.username, err)
	}
}
