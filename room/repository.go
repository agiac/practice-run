package room

import (
	"context"
	"sync"
)

type Repository struct {
	mu   sync.RWMutex
	room map[string]*Room
}

func NewRepository() *Repository {
	return &Repository{
		room: make(map[string]*Room),
	}
}

func (r *Repository) GetRoom(ctx context.Context, id string) (*Room, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if room, ok := r.room[id]; ok {
		return room, nil
	}

	return nil, nil
}

func (r *Repository) UpdateRoom(ctx context.Context, room *Room) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.room[room.Id] = room

	return nil
}
