package user_service

import (
	"context"
	"strings"
	"time"
	"unicode"

	"github.com/suprt/gRPC1/internal/model"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repository UserRepository
}

type UserRepository interface {
	Create(ctx context.Context, name, email, password, role string, createdAt, updatedAt time.Time) (int64, error)
	Get(ctx context.Context, id int64) (model.UserGetResponse, error)
	Update(ctx context.Context, id int64, name, email *string, updatedAt time.Time) error
	Delete(ctx context.Context, id int64) error
}

func New(repository UserRepository) *Service {
	return &Service{repository: repository}
}

func (s *Service) Create(ctx context.Context, name, email, password, passwordConfirm, role string) (int64, error) {
	if strings.ContainsFunc(password, unicode.IsSpace) {
		return 0, model.ErrPasswordSpace
	}
	if password != passwordConfirm {
		return 0, model.ErrPasswordMatch
	}
	if name == "" {
		return 0, model.ErrEmptyName
	}
	if email == "" {
		return 0, model.ErrEmptyEmail
	}
	if password == "" {
		return 0, model.ErrEmptyPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	now := time.Now()
	return s.repository.Create(ctx, name, email, string(hash), role, now, now)
}

func (s *Service) Get(ctx context.Context, id int64) (model.UserGetResponse, error) {
	return s.repository.Get(ctx, id)
}

func (s *Service) Update(ctx context.Context, id int64, name, email *string) error {
	if name != nil && *name == "" {
		return model.ErrEmptyName
	}
	if email != nil && *email == "" {
		return model.ErrEmptyEmail
	}

	return s.repository.Update(ctx, id, name, email, time.Now())
}

func (s *Service) Delete(ctx context.Context, id int64) error {
	return s.repository.Delete(ctx, id)
}
