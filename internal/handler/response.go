package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type ErrorResponse struct {
	Error struct {
		Message string `json:"message"`
		Status  int    `json:"status"`

		Details map[string]string `json:"details,omitempty"`
	} `json:"error"`
}

func JSONError(w http.ResponseWriter, r *http.Request, status int, message string, internalErr error) {
	if internalErr != nil {
		slog.Error("API Error",
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"internal_error", internalErr.Error(),
		)
	}

	resp := ErrorResponse{}
	resp.Error.Message = message
	resp.Error.Status = status

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}

type SuccessResponse struct {
	Data any `json:"data"`
	// TODO: add meta for pagination
}

func JSONResponse(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := SuccessResponse{
		Data: data,
	}

	json.NewEncoder(w).Encode(resp)
}
