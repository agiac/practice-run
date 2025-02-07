package room

import (
	"context"
	"fmt"
	"sync"
)

type Event interface {
	Name() string
}

type Member interface {
	Username() string
	Notify(event Event)
}

type Room struct {
	name string

	mu      sync.Mutex
	members map[string]Member
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) Members() ([]Member, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	members := make([]Member, 0, len(r.members))
	for _, member := range r.members {
		members = append(members, member)
	}

	return members, nil
}

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
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.members[m.Username()]
	if ok {
		return fmt.Errorf("%s already in %s", m.Username(), r.name)
	}

	r.members[m.Username()] = m

	return nil
}

func (s *Service) RemoveMember(ctx context.Context, r *Room, m Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.members, m.Username())

	return nil
}

func (s *Service) SendMessage(ctx context.Context, r *Room, m Member, message string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, member := range r.members {
		if member.Username() == m.Username() {
			continue
		}

		member.Notify(&MessageReceivedEvent{
			RoomName:   r.name,
			SenderName: m.Username(),
			Message:    message,
		})
	}

	return nil
}
