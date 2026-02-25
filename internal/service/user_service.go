package service

import (
	"context"
	"errors"

	"github.com/AfshinJalili/gonod/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type UserService struct {
	repo domain.UserRepository
}

func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, email, password string) error {
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &domain.User{
		Email:    email,
		Password: string(hashedBytes),
	}

	return s.repo.CreateUser(ctx, user)
}

func (s *UserService) Login(ctx context.Context, email, password string) error {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return ErrInvalidCredentials
	}

	return nil
}
