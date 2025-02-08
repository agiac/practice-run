package provider

import (
	"practice-run/chat"
	"practice-run/handler"

	"github.com/gorilla/websocket"
)

func ChatService() *chat.Service {
	return chat.NewService()
}

func WebSocketHandler() *handler.WebSocketHandler {
	return handler.NewWebSocketHandler(
		&websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		ChatService(),
	)
}
