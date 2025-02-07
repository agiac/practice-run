package handler

import (
	"log"
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

func (m *WSChatMember) Notify(event string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	err := m.conn.WriteMessage(websocket.TextMessage, []byte(event))
	if err != nil {
		log.Printf("Error: failed to notify member %s: %v", m.username, err)
	}
}
