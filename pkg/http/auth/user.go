package auth

import (
	"context"

	"github.com/pkg/errors"
)

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
