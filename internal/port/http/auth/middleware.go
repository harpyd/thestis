package auth

import (
	"context"
	"net/http"

	"github.com/harpyd/thestis/internal/port/http/httperr"
)

func Middleware(provider Provider) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, err := provider.AuthenticatedUser(ctx, r)
			if err != nil {
				httperr.Unauthorized("unable-to-verify-user", err, w, r)

				return
			}

			ctx = context.WithValue(ctx, userCtxKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
