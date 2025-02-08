package chat

import (
	"context"
	"fmt"
	"sync"
)

type Service struct {
	mtx   sync.Mutex
	rooms map[string]*Room
}

func NewService() *Service {
	return &Service{
		mtx:   sync.Mutex{},
		rooms: make(map[string]*Room),
	}
}

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

func (r *Service) GetRoom(ctx context.Context, roomName string) (*Room, error) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	room, ok := r.rooms[roomName]
	if !ok {
		return nil, nil
	}

	return room, nil
}

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
