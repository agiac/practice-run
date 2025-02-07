package room

import (
	"context"
	"fmt"
	"slices"
	"sync"
)

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
		return fmt.Errorf("already a room member")
	}

	r.members[m.Username()] = m

	s.broadcastEvent(r, &MemberJoinedEvent{
		RoomName:   r.name,
		MemberName: m.Username(),
	}, m)

	return nil
}

func (s *Service) RemoveMember(ctx context.Context, r *Room, m Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.members, m.Username())

	s.broadcastEvent(r, &MemberLeftEvent{
		RoomName:   r.name,
		MemberName: m.Username(),
	}, m)

	return nil
}

func (s *Service) SendMessage(ctx context.Context, r *Room, m Member, message string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.members[m.Username()]
	if !ok {
		return fmt.Errorf("not a room member")
	}

	s.broadcastEvent(r, &MessageReceivedEvent{
		RoomName:   r.name,
		SenderName: m.Username(),
		Message:    message,
	}, m)

	return nil
}

func (s *Service) broadcastEvent(r *Room, e Event, exclude ...Member) {
	for _, member := range r.members {
		if slices.IndexFunc(exclude, func(i Member) bool {
			return i.Username() == member.Username()
		}) != -1 {
			continue
		}

		member.Notify(e)
	}
}
