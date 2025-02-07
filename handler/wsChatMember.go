package handler

import (
	"fmt"
	"log"
	"practice-run/room"
	"sync"

	"github.com/gorilla/websocket"
)

type WSChatMember struct {
	mu   sync.Mutex
	conn *websocket.Conn

	username string
}

func NewChatMember(username string, conn *websocket.Conn) *WSChatMember {
	return &WSChatMember{username: username, mu: sync.Mutex{}, conn: conn}
}

func (m *WSChatMember) Username() string {
	return m.username
}

func (m *WSChatMember) Notify(event room.Event) {
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

func (m *WSChatMember) notify(message string) {
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to notify member %s: %v", m.username, err)
	}
}

func (m *WSChatMember) WriteMessage(message string) {
	err := m.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		log.Printf("Error: failed to write message to member %s: %v", m.username, err)
	}
}
