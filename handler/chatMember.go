package handler

import (
	"log"
	"practice-run/chat"
	"sync"

	"github.com/gorilla/websocket"
)

type EventHandler interface {
	Handle(event chat.Event, m *ChatMember) error
}

type ChatMember struct {
	mu   sync.Mutex
	conn *websocket.Conn

	username string
	handlers map[string]EventHandler
}

func NewChatMember(username string, conn *websocket.Conn) *ChatMember {
	return &ChatMember{
		username: username,
		conn:     conn,
		handlers: map[string]EventHandler{
			chat.MessageReceivedEventName: &MessageReceivedHandler{},
			chat.MemberJoinedEventName:    &MemberJoinedHandler{},
			chat.MemberLeftEventName:      &MemberLeftHandler{},
		},
	}
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
	handler, ok := m.handlers[event.Name()]
	if !ok {
		log.Printf("Error: failed to notify member %s: unknown event %s", m.username, event.Name())
		return
	}

	err := handler.Handle(event, m)
	if err != nil {
		log.Printf("Error: failed to notify member %s: %v", m.username, err)
	}
}
