package middlewares

import (
	"context"
	"net/http"

	"github.com/varik-08/gw_chat/config"
)

func AppMiddleware(app *config.App) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "app", app)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
