package authprovider

import (
	"context"
	"net/http"
	"strings"

	fireauth "firebase.google.com/go/auth"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/pkg/http/auth"
)

type FirebaseProvider struct {
	authClient *fireauth.Client
}

func Firebase(authClient *fireauth.Client) FirebaseProvider {
	return FirebaseProvider{
		authClient: authClient,
	}
}

var (
	errEmptyBearerToken  = errors.New("empty bearer token")
	errUnableToVerifyJWT = errors.New("unable to verify jwt")
)

func (p FirebaseProvider) AuthenticatedUser(ctx context.Context, r *http.Request) (auth.User, error) {
	bearerToken := tokenFromHeader(r)

	if bearerToken == "" {
		return auth.User{}, errEmptyBearerToken
	}

	token, err := p.authClient.VerifyIDToken(ctx, bearerToken)
	if err != nil {
		return auth.User{}, errUnableToVerifyJWT
	}

	return auth.User{
		UUID:        token.UID,
		DisplayName: token.Claims["name"].(string),
	}, err
}

func tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}
