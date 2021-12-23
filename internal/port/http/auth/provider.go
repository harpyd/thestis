package auth

import (
	"context"
	"net/http"
)

type Provider interface {
	AuthenticatedUser(ctx context.Context, r *http.Request) (User, error)
}
