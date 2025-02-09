package chat

import (
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
