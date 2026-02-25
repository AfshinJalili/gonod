package handler

import (
	"encoding/json"
	"net/http"
)

type HealthHanlder struct{}

func NewHealthHanlder() http.Handler {
	return &HealthHanlder{}
}

func (h *HealthHanlder) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})
}
