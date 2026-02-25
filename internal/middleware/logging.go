package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (sr *statusRecorder) WriteHeader(code int) {
	sr.status = code
	sr.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		sr := &statusRecorder{ResponseWriter: w, status: 200}
		defer func() {
			slog.Info("request",
				"method", r.Method,
				"path", r.URL.Path,
				"status", sr.status,
				"duration", time.Since(start).String(),
				"ip", r.RemoteAddr,
			)
		}()

		next.ServeHTTP(sr, r)

	})
}
