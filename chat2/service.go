package chat2

import (
	"context"
	"fmt"
	"practice-run/room2"
)

//go:generate mockgen -destination mocks/mock_room_service.go -mock_names roomService=RoomService -package=mocks . roomService
type roomService interface {
	CreateRoom(ctx context.Context, roomName string) (*room2.Room, error)
	GetRoom(ctx context.Context, roomName string) (*room2.Room, error)
}

type Service struct {
	rs roomService
}

func NewService(rs roomService) *Service {
	return &Service{
		rs: rs,
	}
}

func (c *Service) AddMemberToRoom(ctx context.Context, roomName string, member room2.Member) error {
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

	err = r.AddMember(ctx, member)
	if err != nil {
		return fmt.Errorf("failed to add member to room: %w", err)
	}

	return nil
}

func (c *Service) RemoveMemberFromRoom(ctx context.Context, roomName string, member room2.Member) error {
	r, err := c.rs.GetRoom(ctx, roomName)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if r == nil {
		return nil
	}

	err = r.RemoveMember(ctx, member)
	if err != nil {
		return fmt.Errorf("failed to remove member from room: %w", err)
	}

	return nil
}

func (c *Service) SendMessageToRoom(ctx context.Context, roomName string, member room2.Member, message string) error {
	r, err := c.rs.GetRoom(ctx, roomName)
	if err != nil {
		return fmt.Errorf("failed to get room: %w", err)
	}

	if r == nil {
		return fmt.Errorf("room not found")
	}

	err = r.SendMessage(ctx, member, message)
	if err != nil {
		return fmt.Errorf("failed to send message to room: %w", err)
	}

	return nil
}
