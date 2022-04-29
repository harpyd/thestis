package rest

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type AuthProvider interface {
	AuthenticateUser(ctx context.Context, r *http.Request) (User, error)
}

type User struct {
	UUID        string
	DisplayName string
}

type ctxKey int

const userCtxKey ctxKey = iota

var errNoUserInCtx = errors.New("no user in context")

func UserFromCtx(ctx context.Context) (User, error) {
	u, ok := ctx.Value(userCtxKey).(User)
	if !ok {
		return User{}, errNoUserInCtx
	}

	return u, nil
}

func AuthMiddleware(provider AuthProvider) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			user, err := provider.AuthenticateUser(ctx, r)
			if err != nil {
				Unauthorized("unable-to-verify-user", err, w, r)

				return
			}

			ctx = context.WithValue(ctx, userCtxKey, user)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
