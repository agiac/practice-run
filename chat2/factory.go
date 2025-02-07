package chat2

import "practice-run/room"

func MakeChatService() *Service {
	return NewService(room.MakeRoomService())
}
