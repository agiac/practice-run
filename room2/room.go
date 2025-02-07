package room2

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

func (r *Room) AddMember(ctx context.Context, member Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.members[member.Username()] = member

	r.broadcastEvent(&MemberJoinedEvent{
		RoomName:   r.name,
		MemberName: member.Username(),
	}, member)

	return nil
}

func (r *Room) RemoveMember(ctx context.Context, member Member) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	delete(r.members, member.Username())

	r.broadcastEvent(&MemberLeftEvent{
		RoomName:   r.name,
		MemberName: member.Username(),
	}, member)

	return nil
}

func (r *Room) SendMessage(ctx context.Context, member Member, message string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.members[member.Username()]
	if !ok {
		return fmt.Errorf("not a room member")
	}

	r.broadcastEvent(&MessageReceivedEvent{
		RoomName:   r.name,
		SenderName: member.Username(),
		Message:    message,
	}, member)

	return nil
}

func (r *Room) broadcastEvent(event Event, exclude ...Member) {
	for _, member := range r.members {
		if slices.IndexFunc(exclude, func(i Member) bool {
			return i.Username() == member.Username()
		}) != -1 {
			continue
		}

		member.Notify(event)
	}
}
