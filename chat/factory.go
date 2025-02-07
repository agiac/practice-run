package chat

import "practice-run/room"

func MakeChatService() *Service {
	return NewService(room.MakeRoomService())
}
