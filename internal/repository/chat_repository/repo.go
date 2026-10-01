package chat_repository

import (
	"context"
	"sync"
	"time"

	"github.com/suprt/gRPC1/internal/model"
)

type Repo struct {
	mu        sync.RWMutex
	chats     map[int64]*chat
	idCounter int64
}

type chat struct {
	usernames []string
	messages  []message
}

type message struct {
	from      string
	text      string
	timestamp time.Time
}

func New() *Repo {
	return &Repo{
		chats: make(map[int64]*chat),
	}
}

func (r *Repo) Create(ctx context.Context, usernames []string) (int64, error) {
	if err := ctx.Err(); err != nil {
		return 0, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if len(usernames) == 0 {
		return 0, model.ErrEmptyUsernames
	}
	r.idCounter++
	id := r.idCounter
	r.chats[id] = &chat{
		usernames: usernames,
		messages:  make([]message, 0),
	}
	return id, nil
}
func (r *Repo) Delete(ctx context.Context, id int64) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.chats[id]; !ok {
		return model.ErrChatNotFound
	}
	delete(r.chats, id)
	return nil
}
func (r *Repo) SendMessage(ctx context.Context, from, text string, id int64, timestamp time.Time) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.chats[id]; !ok {
		return model.ErrChatNotFound
	}
	r.chats[id].messages = append(r.chats[id].messages, message{from: from, text: text, timestamp: timestamp})
	return nil
}
