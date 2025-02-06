package user

import (
	"context"
	"sync"
)

type Repository struct {
	mtx sync.RWMutex
	m   map[string]*User
}

func NewRepository() *Repository {
	return &Repository{
		m: make(map[string]*User),
	}
}

func (r *Repository) UpdateUser(ctx context.Context, u *User) error {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	r.m[u.ID] = u

	return nil
}

func (r *Repository) GetUser(ctx context.Context, id int) (*User, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if u, ok := r.m[id]; ok {
		return u, nil
	}

	return nil, nil
}
