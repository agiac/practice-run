package chat

import (
	"context"
	"fmt"
	"practice-run/internal/room"
)

//go:generate mockgen -destination mocks/mock_room_service.go -mock_names roomRepository=RoomRepository -package=mocks . roomRepository
type roomRepository interface {
	CreateRoom(ctx context.Context, roomName string) (*room.Room, error)
	GetRoom(ctx context.Context, roomName string) (*room.Room, error)
}

//go:generate mockgen -destination mocks/mock_room_manager.go -mock_names roomManager=RoomManager -package=mocks . roomManager
type roomManager interface {
	AddMember(ctx context.Context, r *room.Room, m room.Member) error
	RemoveMember(ctx context.Context, r *room.Room, m room.Member) error
	SendMessage(ctx context.Context, r *room.Room, m room.Member, message string) error
}

type Service struct {
	rs roomRepository
	rm roomManager
}

func NewService(rs roomRepository, rm roomManager) *Service {
	return &Service{
		rs: rs,
		rm: rm,
	}
}

func (c *Service) AddMemberToRoom(ctx context.Context, roomName string, member room.Member) error {
	r, err := c.rs.GetRoom(ctx, roomName)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if r == nil {
		r, err = c.rs.CreateRoom(ctx, roomName)
		if err != nil {
			return fmt.Errorf("failed to create room: %w", err)
		}
	}

	err = c.rm.AddMember(ctx, r, member)
	if err != nil {
		return fmt.Errorf("failed to add member to room: %w", err)
	}

	return nil
}

func (c *Service) RemoveMemberFromRoom(ctx context.Context, roomName string, member room.Member) error {
	r, err := c.rs.GetRoom(ctx, roomName)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if r == nil {
		return nil
	}

	err = c.rm.RemoveMember(ctx, r, member)
	if err != nil {
		return fmt.Errorf("failed to remove member from room: %w", err)
	}

	return nil
}

func (c *Service) SendMessageToRoom(ctx context.Context, roomName string, member room.Member, message string) error {
	r, err := c.rs.GetRoom(ctx, roomName)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if r == nil {
		return fmt.Errorf("room not found")
	}

	err = c.rm.SendMessage(ctx, r, member, message)
	if err != nil {
		return fmt.Errorf("failed to send message to room: %w", err)
	}

	return nil
}
