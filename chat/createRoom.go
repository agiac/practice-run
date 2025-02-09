package chat

import (
	"context"
	"fmt"
)

func (r *Service) CreateRoom(ctx context.Context, name string) (*Room, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.rooms[name]
	if ok {
		return nil, fmt.Errorf("room already exists")
	}

	room = &Room{
		name:    name,
		members: make(map[string]Member),
	}

	r.rooms[name] = room

	return room, nil
}
