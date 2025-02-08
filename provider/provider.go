package provider

import (
	"net/http"
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
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		ChatService(),
	)
}
