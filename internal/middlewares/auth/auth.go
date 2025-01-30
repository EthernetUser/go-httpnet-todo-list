package auth

import (
	"context"
	"go-httpnet-todo-list/internal/consts"
	"net/http"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// TODO: Add middleware to check if user is authenticated
		ctx := context.WithValue(r.Context(), consts.AuthUserIdKey, 1)
		req := r.WithContext(ctx)
		next.ServeHTTP(w, req)
	})
}