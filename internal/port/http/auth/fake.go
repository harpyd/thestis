package auth

import (
	"context"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/golang-jwt/jwt/request"

	"github.com/harpyd/thestis/pkg/httperr"
)

func FakeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			httperr.BadRequest("unable-to-get-jwt", err, w, r)

			return
		}

		if !token.Valid {
			httperr.BadRequest("invalid-jwt", nil, w, r)

			return
		}

		ctx := context.WithValue(r.Context(), userCtxKey, User{
			UUID:        claims["userUuid"].(string),
			DisplayName: claims["name"].(string),
		})

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
