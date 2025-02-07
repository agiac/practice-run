package provider

import (
	"practice-run/internal/chat"
	"practice-run/internal/handler"
	room2 "practice-run/internal/room"

	"github.com/gorilla/websocket"
)

func RoomRepository() *room2.Repository {
	return room2.NewRepository()
}

func RoomManager() *room2.Manager {
	return room2.NewManager()
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
