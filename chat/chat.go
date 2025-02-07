package chat

import (
	"context"
	"fmt"
	"practice-run/room"
	"sync"
)

//go:generate mockgen -destination mocks/mock_room_service.go -package=mocks . RoomService
type RoomService interface {
	CreateRoom(ctx context.Context, roomName string) (*room.Room, error)
}

type Service struct {
	mu    sync.Mutex
	rooms map[string]*room.Room

	rs RoomService
}

func NewService(rs RoomService) *Service {
	return &Service{
		mu:    sync.Mutex{},
		rooms: make(map[string]*room.Room),
		rs:    rs,
	}
}

func (c *Service) AddMemberToRoom(ctx context.Context, roomName string, member room.Member) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var err error
	r, ok := c.rooms[roomName]
	if !ok {
		r, err = c.rs.CreateRoom(ctx, roomName)
		if err != nil {
			return fmt.Errorf("failed to create room: %w", err)
		}

		c.rooms[roomName] = r
	}

	err = r.AddMember(ctx, member)
	if err != nil {
		return fmt.Errorf("failed to add member to room: %w", err)
	}

	return nil
}

func (c *Service) RemoveMemberFromRoom(ctx context.Context, roomName string, member room.Member) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	r, ok := c.rooms[roomName]
	if !ok {
		return nil
	}

	err := r.RemoveMember(ctx, member)
	if err != nil {
		return fmt.Errorf("failed to remove member from room: %w", err)
	}

	return nil
}

func (c *Service) SendMessageToRoom(ctx context.Context, roomName string, member room.Member, message string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	r, ok := c.rooms[roomName]
	if !ok {
		return fmt.Errorf("room %s not found", roomName)
	}

	err := r.SendMessage(ctx, member, message)
	if err != nil {
		return fmt.Errorf("failed to send message to room: %w", err)
	}

	return nil
}
