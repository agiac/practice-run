package chat

import (
	"context"
	"fmt"
)

func (r *Service) SendMessage(ctx context.Context, roomName string, member Member, message string) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.rooms[roomName]
	if !ok {
		return fmt.Errorf("room not found")
	}

	err := room.sendMessage(member, message)
	if err != nil {
		return fmt.Errorf("failed to send message to room: %w", err)
	}

	return nil
}
