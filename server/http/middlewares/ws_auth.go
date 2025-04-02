package middlewares

import (
	"context"
	"net/http"

	"github.com/varik-08/gw_chat/internal/pkg/jwt"
)

func WsAuthMiddleware(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenStr := r.URL.Query().Get("token")
			if tokenStr == "" {
				http.Error(w, "Missing token in URL", http.StatusUnauthorized)
				return
			}

			claims, err := jwt.ValidateToken(secretKey, tokenStr)
			if err != nil {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
			ctx = context.WithValue(ctx, UsernameKey, claims.Username)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
