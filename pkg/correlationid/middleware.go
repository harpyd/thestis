package correlationid

import (
	"net/http"

	"github.com/google/uuid"
)

func Middleware(header string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			correlationID := r.Header.Get(header)

			if correlationID == "" {
				correlationID = uuid.New().String()
			}

			ctx = AssignToCtx(ctx, correlationID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
