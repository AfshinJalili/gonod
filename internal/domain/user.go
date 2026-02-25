package domain

import (
	"context"
	"time"
)

type User struct {
	ID string
	Email string
	Passowrd string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserRepositry interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}