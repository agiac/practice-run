package provider

import (
	"practice-run/internal/chat"
	"practice-run/internal/handler"

	"github.com/gorilla/websocket"
)

func ChatService() *chat.Service {
	return chat.NewService()
}

func WebSocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}

func WebSocketHandler() *handler.WebSocketHandler {
	return handler.NewWebSocketHandler(
		WebSocketUpgrader(),
		ChatService(),
	)
}
