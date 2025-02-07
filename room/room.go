package room

import "sync"

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
