package chat

import (
	"context"
	"fmt"
	"slices"
)

type Member interface {
	Username() string
	Notify(event Event)
}

type Room struct {
	name string

	members map[string]Member // protected from concurrent access by the service layer
}

func (r *Room) Name() string {
	return r.name
}

func (r *Room) getMembers(ctx context.Context) ([]Member, error) {
	members := make([]Member, 0, len(r.members))
	for _, member := range r.members {
		members = append(members, member)
	}

	return members, nil
}

func (r *Room) addMember(ctx context.Context, member Member) error {
	_, ok := r.members[member.Username()]
	if ok {
		return fmt.Errorf("member already exists")
	}

	r.members[member.Username()] = member

	r.broadcastEvent(&MemberJoinedEvent{
		RoomName:   r.Name(),
		MemberName: member.Username(),
	}, member)

	return nil
}

func (r *Room) removeMember(ctx context.Context, member Member) error {
	if _, ok := r.members[member.Username()]; !ok {
		return fmt.Errorf("not a room member")
	}

	delete(r.members, member.Username())

	r.broadcastEvent(&MemberLeftEvent{
		RoomName:   r.Name(),
		MemberName: member.Username(),
	}, member)

	return nil
}

func (r *Room) sendMessage(ctx context.Context, member Member, message string) error {
	_, ok := r.members[member.Username()]
	if !ok {
		return fmt.Errorf("not a room member")
	}

	r.broadcastEvent(&MessageReceivedEvent{
		RoomName:   r.Name(),
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
