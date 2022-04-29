package firebase

import (
	"context"
	stdhttp "net/http"
	"strings"

	"firebase.google.com/go/auth"
	"github.com/pkg/errors"

	"github.com/harpyd/thestis/internal/core/interface/rest"
)

type Provider struct {
	authClient *auth.Client
}

func NewProvider(authClient *auth.Client) Provider {
	return Provider{
		authClient: authClient,
	}
}

var (
	errEmptyBearerToken  = errors.New("empty bearer token")
	errUnableToVerifyJWT = errors.New("unable to verify jwt")
	errInvalidClaimType  = errors.New("invalid claim type")
)

func (p Provider) AuthenticateUser(ctx context.Context, r *stdhttp.Request) (rest.User, error) {
	bearerToken := tokenFromHeader(r)

	if bearerToken == "" {
		return rest.User{}, errEmptyBearerToken
	}

	token, err := p.authClient.VerifyIDToken(ctx, bearerToken)
	if err != nil {
		return rest.User{}, errUnableToVerifyJWT
	}

	displayName, ok := token.Claims["name"].(string)
	if !ok {
		return rest.User{}, errInvalidClaimType
	}

	return rest.User{
		UUID:        token.UID,
		DisplayName: displayName,
	}, nil
}

func tokenFromHeader(r *stdhttp.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}
