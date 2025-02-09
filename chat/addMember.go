package chat

import (
	"context"
	"fmt"
)

func (r *Service) AddMember(ctx context.Context, roomName string, member Member) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.rooms[roomName]
	if !ok {
		return fmt.Errorf("room not found")
	}

	err := room.addMember(member)
	if err != nil {
		return fmt.Errorf("failed to add member to room: %w", err)
	}

	return nil
}
