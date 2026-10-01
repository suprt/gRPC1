package model

import (
	"errors"
	"time"
)

type UserGetResponse struct {
	ID      int64
	Name    string
	Email   string
	Role    string
	Created time.Time
	Updated time.Time
}

var (
	ErrEmptyName      = errors.New("empty name")
	ErrEmptyEmail     = errors.New("empty email")
	ErrPasswordSpace  = errors.New("password must not contain spaces")
	ErrPasswordMatch  = errors.New("password must match")
	ErrUserNotFound   = errors.New("user not found")
	ErrEmptyPassword  = errors.New("empty password")
	ErrChatNotFound   = errors.New("chat not found")
	ErrEmptyUsernames = errors.New("empty usernames")
)
