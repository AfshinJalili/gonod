package repository

import (
	"context"
	"database/sql"

	"github.com/AfshinJalili/gonod/internal/domain"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepositry {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, u *domain.User) error {	
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`
	err := r.db.QueryRowContext(ctx, query, u.Email, u.Passowrd).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)

	return err
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1`

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.Email, &u.Passowrd, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
