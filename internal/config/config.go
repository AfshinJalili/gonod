package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Environment string
	Port  string
	DBURL string
}

func Load() *Config {
	// Try to load .env file for local development.
	// We ignore the error because in production, we rely on actual system env vars.
	_ = godotenv.Load()

	cfg := &Config{
		Environment: os.Getenv("ENVIRONMENT"),
		Port:  os.Getenv("PORT"),
		DBURL: os.Getenv("DB_URL"),

	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.Environment == "" {
		cfg.Environment = "development"
	}

	return cfg
}
