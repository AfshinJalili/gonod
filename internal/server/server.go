package server

import (
	"net/http"

	"github.com/AfshinJalili/gonod/internal/middleware"
)

func New() http.Handler {
	mux := http.NewServeMux()

	registerRoutes(mux)

	return middleware.Chain(mux, middleware.Logging, middleware.Recover)
}
