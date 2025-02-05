package middlewares

import (
	"fmt"
	"go-httpnet-todo-list/internal/consts"
	"net/http"

	"github.com/google/uuid"
)

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		xRequestIDHeader := r.Header.Get(consts.RequestIdHeader)
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
		w.Header().Set(consts.RequestIdHeader, requestId)
		next.ServeHTTP(w, r)
	})
}
