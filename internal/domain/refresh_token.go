package domain

import (
	"context"
	"time"
)

type RefreshToken struct {
	ID                string
	SessionID         string
	TokenHash         string
	ReplacedByTokenID string
	ExpiresAt         time.Time
	RotatedAt         time.Time
	RevokedAt         time.Time
	CreatedAt         time.Time
}


type RefreshTokenRepository interface {
	CreateRefreshToken(ctx context.Context, refreshToken RefreshToken) error
	GetRefreshTokenByHash(ctx context.Context, tokenHash string) (RefreshToken, error)
	RotateRefreshToken(ctx context.Context, ID string) (RefreshToken, error)
	RevokeRefreshToken(ctx context.Context, ID string) error
}
