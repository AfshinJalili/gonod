package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHandler struct{}

func NewHealthHandler() http.Handler {
	return &HealthHandler{}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}
