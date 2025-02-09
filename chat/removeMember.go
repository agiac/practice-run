package chat

import (
	"context"
	"fmt"
)

func (r *Service) RemoveMember(ctx context.Context, roomName string, member Member) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.rooms[roomName]
	if !ok {
		return fmt.Errorf("room not found")
	}

	err := room.removeMember(member)
	if err != nil {
		return fmt.Errorf("failed to remove member from room: %w", err)
	}

	return nil
}
