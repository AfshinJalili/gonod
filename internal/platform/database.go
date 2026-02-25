package platform

import (
	"database/sql"
	"errors"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib" // The pgx driver
)

func SetupDatabase(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	slog.Info("Successfully connected and pinged PostgreSQL")

	err = runMigrations(dbURL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(dbURL string) error {
	m, err := migrate.New("file://migrations", dbURL)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("Database migrations are already up to date")
			return nil
		}
		return err
	}

	slog.Info("Database migrations applied successfully")
	return nil
}
