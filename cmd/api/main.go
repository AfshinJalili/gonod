package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/AfshinJalili/gonod/internal/config"
	"github.com/AfshinJalili/gonod/internal/platform"
	"github.com/AfshinJalili/gonod/internal/server"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	cfg := config.Load()

	db, err := platform.SetupDatabase(cfg.DBURL)
	if err != nil {
		slog.Error("Failed to setup database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	srv := server.New()

	slog.Info("Starting server", "port", cfg.Port, "environment", cfg.Environment)

	err = http.ListenAndServe(":"+cfg.Port, srv)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
