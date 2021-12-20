package auth

import (
	"context"
	"net/http"
	"strings"

	"firebase.google.com/go/auth"

	"github.com/harpyd/thestis/pkg/httperr"
)

func FirebaseMiddleware(auth *auth.Client) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			bearerToken := tokenFromHeader(r)
			if bearerToken == "" {
				httperr.Unauthorized("empty-bearer-token", nil, w, r)

				return
			}

			token, err := auth.VerifyIDToken(ctx, bearerToken)
			if err != nil {
				httperr.Unauthorized("unable-to-verify-jwt", err, w, r)

				return
			}

			ctx = context.WithValue(ctx, userCtxKey, User{
				UUID:        token.UID,
				DisplayName: token.Claims["name"].(string),
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func tokenFromHeader(r *http.Request) string {
	headerValue := r.Header.Get("Authorization")

	if len(headerValue) > 7 && strings.ToLower(headerValue[0:6]) == "bearer" {
		return headerValue[7:]
	}

	return ""
}
