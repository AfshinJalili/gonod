package domain

import (
	"context"
	"errors"
	"time"
)

type Session struct {
	ID        string
	UserID    string
	UserAgent string
	IP        string
	ExpiresAt time.Time
	RevokedAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
}

type SessionRepository interface {
	CreateSession(ctx context.Context, session *Session) error
	RevokeSession(ctx context.Context, ID string) error
}

var SessionNotFoundErr = errors.New("session not found")