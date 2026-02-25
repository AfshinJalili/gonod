package handler

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/AfshinJalili/gonod/internal/domain"
	"github.com/AfshinJalili/gonod/internal/service"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(us *service.UserService) *AuthHandler {
	return &AuthHandler{userService: us}
}

type AuthRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}

		slog.Error("Failed to register user", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login Successful"})
}
