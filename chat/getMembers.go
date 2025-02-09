package chat

import (
	"context"
	"fmt"
)

func (r *Service) GetMembers(ctx context.Context, roomName string) ([]Member, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.rooms[roomName]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}

	members, err := room.getMembers()
	if err != nil {
		return nil, fmt.Errorf("failed to get members: %w", err)
	}

	return members, nil
}
