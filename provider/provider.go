package provider

import (
	"practice-run/chat"
	"practice-run/handler"
	"practice-run/room"

	"github.com/gorilla/websocket"
)

func RoomService() *room.Service {
	return room.NewService()
}

func ChatService() *chat.Service {
	return chat.NewService(RoomService())
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
