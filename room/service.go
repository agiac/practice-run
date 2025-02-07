package room

import (
	"context"
	"sync"
)

type Service struct {
}

func NewService() *Service {
	return &Service{}
}

func (s *Service) CreateRoom(ctx context.Context, name string) (*Room, error) {
	return &Room{
		name:    name,
		mu:      sync.Mutex{},
		members: make(map[string]Member),
	}, nil
}

func (s *Service) AddMember(ctx context.Context, r *Room, m Member) error {
	return r.addMember(ctx, m)
}

func (s *Service) RemoveMember(ctx context.Context, r *Room, m Member) error {
	return r.removeMember(ctx, m)
}

func (s *Service) SendMessage(ctx context.Context, r *Room, m Member, message string) error {
	return r.sendMessage(ctx, m, message)
}
