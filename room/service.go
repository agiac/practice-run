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
