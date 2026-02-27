package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/AfshinJalili/gonod/internal/domain"
	"github.com/AfshinJalili/gonod/internal/service"
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)

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

func (req *AuthRequest) Sanitize() {
	req.Email = strings.ToLower(strings.TrimSpace(req.Email))
	req.Password = strings.TrimSpace(req.Password)
}

func (req *AuthRequest) Validate() map[string]string {
	errors := make(map[string]string)

	if req.Email == "" {
		errors["email"] = "email is required"
	} else if !emailRegex.MatchString(req.Email) {
		errors["email"] = "must be a valid email format"
	}

	if len(strings.TrimSpace(req.Password)) < 8 {
		errors["password"] = "password must be at least 8 characters long"
	}

	return errors
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, r, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	req.Sanitize()

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		JSONError(w, r, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	err := h.userService.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrDuplicateEmail) {
			JSONError(w, r, http.StatusConflict, "Email already exist", err)
			return
		}

		JSONError(w, r, http.StatusInternalServerError, "Internal server error", err)
		return
	}

	JSONResponse(w, http.StatusCreated, map[string]string{"message": "User registered successfully"})
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req AuthRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		JSONError(w, r, http.StatusBadRequest, "Invalid request payload", err)
		return
	}

	req.Sanitize()

	if validationErrors := req.Validate(); len(validationErrors) > 0 {
		JSONError(w, r, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		JSONError(w, r, http.StatusUnauthorized, "Invalid credentials", err)
		return
	}

	JSONResponse(w, http.StatusOK, map[string]string{"message": "Login Successful"})
}
