package provider

import (
	"practice-run/chat"
	"practice-run/handler"
	"practice-run/room"

	"github.com/gorilla/websocket"
)

func RoomRepository() *room.Repository {
	return room.NewRepository()
}

func RoomManager() *room.Manager {
	return room.NewManager()
}

func ChatService() *chat.Service {
	return chat.NewService(RoomRepository(), RoomManager())
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
