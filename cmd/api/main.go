package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/AfshinJalili/gonod/internal/config"
	"github.com/AfshinJalili/gonod/internal/handler"
	"github.com/AfshinJalili/gonod/internal/platform"
	"github.com/AfshinJalili/gonod/internal/repository"
	"github.com/AfshinJalili/gonod/internal/server"
	"github.com/AfshinJalili/gonod/internal/service"
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

	userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
	authHandler := handler.NewAuthHandler(userService)

	srv := server.New(authHandler)

	slog.Info("Starting server", "port", cfg.Port, "environment", cfg.Environment)

	err = http.ListenAndServe(":"+cfg.Port, srv)
	if err != nil {
		slog.Error("Server failed to start", "error", err)
	}
}
