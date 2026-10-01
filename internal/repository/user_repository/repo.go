package user_repository

import (
	"context"
	"sync"
	"time"

	"github.com/suprt/gRPC1/internal/model"
)

type user struct {
	name      string
	email     string
	password  string
	role      string
	createdAt time.Time
	updatedAt time.Time
}

type Repo struct {
	mu        sync.RWMutex
	users     map[int64]*user
	idCounter int64
}

func New() *Repo {
	return &Repo{
		users: make(map[int64]*user),
	}
}

func (r *Repo) Create(ctx context.Context, name, email, password, role string, createdAt, updatedAt time.Time) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idCounter++
	id := r.idCounter
	r.users[id] = &user{
		name:      name,
		email:     email,
		password:  password,
		role:      role,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	return id, nil
}

func (r *Repo) Get(ctx context.Context, id int64) (model.UserGetResponse, error) {
	if err := ctx.Err(); err != nil {
		return model.UserGetResponse{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	u, ok := r.users[id]
	if !ok {
		return model.UserGetResponse{}, model.ErrUserNotFound
	}
	return model.UserGetResponse{
		ID:      id,
		Name:    u.name,
		Email:   u.email,
		Role:    u.role,
		Created: u.createdAt,
		Updated: u.updatedAt,
	}, nil
}

func (r *Repo) Update(ctx context.Context, id int64, name, email *string, updatedAt time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	u, ok := r.users[id]
	if !ok {
		return model.ErrUserNotFound
	}
	if name != nil {
		u.name = *name
	}
	if email != nil {
		u.email = *email
	}
	u.updatedAt = updatedAt
	return nil
}

func (r *Repo) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	_, ok := r.users[id]
	if !ok {
		return model.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}
