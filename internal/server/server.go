package server

import (
	"net/http"

	"github.com/AfshinJalili/gonod/internal/handler"
	"github.com/AfshinJalili/gonod/internal/middleware"
)

type Server struct {
	mux           *http.ServeMux
	healthHandler http.Handler
	authHandler   *handler.AuthHandler
}

func New(authHandler *handler.AuthHandler) http.Handler {
	s := &Server{
		mux:           http.NewServeMux(),
		healthHandler: handler.NewHealthHandler(),
		authHandler:   authHandler,
	}

	s.registerRoutes()

	return middleware.Chain(s.mux, middleware.Logging, middleware.Recover)
}
