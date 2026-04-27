package repository

import (
	"context"
	"database/sql"
	"sort"

	"github.com/AfshinJalili/gonod/internal/domain"
)


type PostgresSessionRepository struct {
	db *sql.DB
}

func newSessionRepository(db *sql.DB) domain.SessionRepository {
	return &PostgresSessionRepository{db}
}

func (r *PostgresSessionRepository) CreateSession(ctx context.Context, s *domain.Session) error {
	query := `
		INSERT INTO sessions (user_id, user_agent, ip, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRowContext(ctx, query, s.UserID, s.UserAgent, s.ExpiresAt).
		Scan(&s.ID, &s.CreatedAt, &s.UpdatedAt)

	return err
}

func (r *PostgresSessionRepository) RevokeSession(ctx context.Context, ID string) error {
	query := `
		UPDATE sessions
		SET revoked_at = NOW()
		WHERE id = $1
	`
	result, err := r.db.ExecContext(ctx, query, ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return domain.SessionNotFoundErr
	}

	return nil
}