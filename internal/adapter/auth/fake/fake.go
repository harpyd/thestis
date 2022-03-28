package fake

import (
	"context"
	stdhttp "net/http"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/port/http"
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

var (
	errUnableToGetJWT = errors.New("unable to get jwt")
	errInvalidJWT     = errors.New("invalid jwt")
)

func (p Provider) AuthenticateUser(_ context.Context, r *stdhttp.Request) (http.User, error) {
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
		return http.User{}, errUnableToGetJWT
	}

	if !token.Valid {
		return http.User{}, errInvalidJWT
	}

	return http.User{
		UUID:        claims["userUuid"].(string),
		DisplayName: claims["name"].(string),
	}, nil
}
