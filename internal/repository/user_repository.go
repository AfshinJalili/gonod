package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/AfshinJalili/gonod/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &PostgresUserRepository{db: db}
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, u *domain.User) error {
	query := `
		INSERT INTO users (email, password)
		VALUES ($1, $2)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, u.Email, u.Password).
		Scan(&u.ID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			// 23505 is the official Postgres code for unique_violation
			if pgErr.Code == "23505" {
				return domain.ErrDuplicateEmail
			}
		}

		return err
	}

	return nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	query := `
		SELECT id, email, password, created_at, updated_at
		FROM users
		WHERE email = $1`

	var u domain.User
	err := r.db.QueryRowContext(ctx, query, email).
		Scan(&u.ID, &u.Email, &u.Password, &u.CreatedAt, &u.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &u, nil
}
