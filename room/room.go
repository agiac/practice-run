package room

import (
	"context"
	"fmt"
	"sync"
)

type Member interface {
	Username() string
	Notify(event string)
}

type Room struct {
	name string

	mu      sync.Mutex
	members map[string]Member
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) AddMember(ctx context.Context, m Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.members[m.Username()]
	if ok {
		return fmt.Errorf("%s already in %s", m.Username(), r.name)
	}

	r.members[m.Username()] = m

	return nil
}

func (r *Room) RemoveMember(ctx context.Context, m Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.members, m.Username())

	return nil
}

func (r *Room) SendMessage(ctx context.Context, m Member, message string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, member := range r.members {
		if member.Username() == m.Username() {
			continue
		}

		member.Notify(fmt.Sprintf("%s: @%s: %s", r.name, m.Username(), message))
	}

	return nil
}
