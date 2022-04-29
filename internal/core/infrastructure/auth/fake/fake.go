package fake

import (
	"context"
	stdhttp "net/http"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/interface/http"
)

type Provider struct{}

func NewProvider() Provider {
	return Provider{}
}

var (
	errUnableToGetJWT   = errors.New("unable to get jwt")
	errInvalidJWT       = errors.New("invalid jwt")
	errInvalidClaimType = errors.New("invalid claim type")
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

	uuid, ok := claims["userUuid"].(string)
	if !ok {
		return http.User{}, errInvalidClaimType
	}

	displayName, ok := claims["name"].(string)
	if !ok {
		return http.User{}, errInvalidClaimType
	}

	return http.User{
		UUID:        uuid,
		DisplayName: displayName,
	}, nil
}
