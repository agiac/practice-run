package handler

import (
	"practice-run/chat"

	"github.com/gorilla/websocket"
)

func MakeHandler() *WebSocketHandler {
	return NewWebSocketHandler(
		&websocket.Upgrader{ // TODO: upgrader settings
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		chat.MakeChatService(),
	)
}
