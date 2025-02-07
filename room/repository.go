package room

import (
	"context"
	"fmt"
	"sync"
)

type Repository struct {
	mtx   sync.Mutex
	rooms map[string]*Room
}

func NewRepository() *Repository {
	return &Repository{
		mtx:   sync.Mutex{},
		rooms: make(map[string]*Room),
	}
}

func (r *Repository) CreateRoom(ctx context.Context, name string) (*Room, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.rooms[name]
	if ok {
		return nil, fmt.Errorf("room already exists")
	}

	room = &Room{name: name, members: make(map[string]Member)}
	r.rooms[name] = room

	return room, nil
}

func (r *Repository) GetRoom(ctx context.Context, roomName string) (*Room, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.rooms[roomName]
	if !ok {
		return nil, nil
	}

	return room, nil
}
