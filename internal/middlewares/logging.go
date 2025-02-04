package middlewares

import (
	"log/slog"
	"net/http"
	"os"
	"time"
)

type wrappedWriter struct {
	http.ResponseWriter
	statusCode int
}

type Logger interface {
	Info(message string, args ...any)
}

func (w *wrappedWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
	w.statusCode = statusCode
}

func Logging(logger Logger) func(http.Handler) http.Handler {
	if logger == nil {
		logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
	}

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			wrapped := &wrappedWriter{
				ResponseWriter: w,
			}
			next.ServeHTTP(wrapped, r)
			requestId := w.Header().Get("X-Request-ID")
			logger.Info(
				"Request completed",
				"method",
				r.Method,
				"path",
				r.URL.Path,
				"request_id",
				requestId,
				"status_code",
				wrapped.statusCode,
				"duration",
				time.Since(start).String(),
			)
		})
	}
}
