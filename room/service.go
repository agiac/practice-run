package room

import (
	"context"
	"fmt"
	"slices"

	"github.com/google/uuid"
)

type repository interface {
	UpdateRoom(ctx context.Context, r *Room) error
	GetRoom(ctx context.Context, id string) (*Room, error)
}

type Service struct {
	repository repository
}

func NewService(repository repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) UserCreatesRoom(ctx context.Context, userId string, roomName string) (*Room, error) {
	roomId, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("failed to generate room id: %w", err)
	}

	room := &Room{
		Id:      roomId.String(),
		Name:    roomName,
		Members: []string{userId},
	}

	if err = s.repository.UpdateRoom(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to create room: %w", err)
	}

	return room, nil
}

func (s *Service) UserJoinsRoom(ctx context.Context, userId string, roomId string) (*Room, error) {
	room, err := s.repository.GetRoom(ctx, roomId)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	room.Members = append(room.Members, userId)

	if err = s.repository.UpdateRoom(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	return room, nil
}

func (s *Service) UserLeavesRoom(ctx context.Context, userId string, roomId string) (*Room, error) {
	room, err := s.repository.GetRoom(ctx, roomId)
	if err != nil {
		return nil, fmt.Errorf("failed to get room: %w", err)
	}

	idx := slices.Index(room.Members, userId)
	if idx == -1 {
		return nil, fmt.Errorf("user is not a member of the room")
	}

	room.Members = append(room.Members[:idx], room.Members[idx+1:]...)

	if err = s.repository.UpdateRoom(ctx, room); err != nil {
		return nil, fmt.Errorf("failed to update room: %w", err)
	}

	return room, nil
}
