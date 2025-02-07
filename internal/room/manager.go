package room

import "context"

type Manager struct {
}

func NewManager() *Manager {
	return &Manager{}
}

func (m2 *Manager) GetMembers(ctx context.Context, r *Room) ([]Member, error) {
	return r.getMembers(ctx)
}

func (m2 *Manager) AddMember(ctx context.Context, r *Room, m Member) error {
	return r.addMember(ctx, m)
}

func (m2 *Manager) RemoveMember(ctx context.Context, r *Room, m Member) error {
	return r.removeMember(ctx, m)
}

func (m2 *Manager) SendMessage(ctx context.Context, r *Room, m Member, message string) error {
	return r.sendMessage(ctx, m, message)
}
