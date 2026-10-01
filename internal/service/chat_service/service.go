package chat_service

import (
	"context"
	"time"
)

type Service struct {
	repository ChatRepository
}

type ChatRepository interface {
	Create(ctx context.Context, usernames []string) (int64, error)
	Delete(ctx context.Context, id int64) error
	SendMessage(ctx context.Context, from, text string, id int64, timestamp time.Time) error
}

func New(repository ChatRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, usernames []string) (int64, error) {
	return s.repository.Create(ctx, usernames)
}
func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
func (s *Service) SendMessage(ctx context.Context, from, text string, id int64, timestamp time.Time) error {
	return s.repository.SendMessage(ctx, from, text, id, timestamp)
}
