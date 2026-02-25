package server

import (
	"net/http"

	"github.com/AfshinJalili/gonod/internal/handler"
)

func registerRoutes(mux *http.ServeMux) {
	
	healthHanlder := handler.NewHealthHanlder()

	mux.Handle("/health", healthHanlder)
}