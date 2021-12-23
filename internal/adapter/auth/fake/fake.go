package fake

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/pkg/http/auth"
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

var (
	errUnableToGetJWT = errors.New("unable to get jwt")
	errInvalidJWT     = errors.New("invalid jwt")
)

func (p Provider) AuthenticatedUser(_ context.Context, r *http.Request) (auth.User, error) {
	var claims jwt.MapClaims

	token, err := request.ParseFromRequest(
		r,
		request.AuthorizationHeaderExtractor,
		func(token *jwt.Token) (i interface{}, e error) {
			return []byte("mock-secret"), nil
		},
		request.WithClaims(&claims),
	)
	if err != nil {
		return auth.User{}, errUnableToGetJWT
	}

	if !token.Valid {
		return auth.User{}, errInvalidJWT
	}

	return auth.User{
		UUID:        claims["userUuid"].(string),
		DisplayName: claims["name"].(string),
	}, nil
}
