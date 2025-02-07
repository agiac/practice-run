package room

import (
	"context"
	"errors"
)

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m2 *Manager) GetMembers(ctx context.Context, r *Room) ([]Member, error) {
	if r == nil {
		return nil, errors.New("room cannot be nil")
	}

	return r.getMembers(ctx)
}

func (m2 *Manager) AddMember(ctx context.Context, r *Room, m Member) error {
	if r == nil {
		return errors.New("room cannot be nil")
	}

	return r.addMember(ctx, m)
}

func (m2 *Manager) RemoveMember(ctx context.Context, r *Room, m Member) error {
	if r == nil {
		return errors.New("room cannot be nil")
	}

	return r.removeMember(ctx, m)
}

func (m2 *Manager) SendMessage(ctx context.Context, r *Room, m Member, message string) error {
	if r == nil {
		return errors.New("room cannot be nil")
	}

	return r.sendMessage(ctx, m, message)
}
