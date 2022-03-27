package auth

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
)

type Provider interface {
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
