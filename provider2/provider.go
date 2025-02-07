package provider2

import (
	"practice-run/chat2"
	"practice-run/handler2"
	"practice-run/room2"

	"github.com/gorilla/websocket"
)

func RoomRepository() *room2.Repository {
	return room2.NewRepository()
}

func RoomManager() *room2.Manager {
	return room2.NewManager()
}

func ChatService() *chat2.Service {
	return chat2.NewService(RoomRepository(), RoomManager())
}

func WebSocketUpgrader() *websocket.Upgrader {
	return &websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
}

func WebSocketHandler() *handler2.WebSocketHandler {
	return handler2.NewWebSocketHandler(
		WebSocketUpgrader(),
		ChatService(),
	)
}
