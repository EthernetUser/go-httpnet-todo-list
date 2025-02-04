package middlewares

import (
	"fmt"
	"net/http"

	"github.com/google/uuid"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestIDHeader := r.Header.Get("X-Request-ID")
		if xRequestIDHeader != "" {
			next.ServeHTTP(w, r)
			return
		}

		id := uuid.New().String()
		requestId := fmt.Sprintf(
			"%s-%s-%s-%s",
			r.Host,
			r.RemoteAddr,
			r.UserAgent(),
			id,
		)
		w.Header().Set("X-Request-ID", requestId)
		next.ServeHTTP(w, r)
	})
}
