package v1

import (
	"net/http"

	"github.com/harpyd/thestis/internal/port/http/auth"
	"github.com/harpyd/thestis/internal/port/http/httperr"
)

func unmarshalUser(w http.ResponseWriter, r *http.Request) (auth.User, bool) {
	user, err := auth.UserFromCtx(r.Context())
	if err != nil {
		httperr.Unauthorized(string(ErrorSlugUnauthorizedUser), err, w, r)

		return auth.User{}, false
	}

	return user, true
}
