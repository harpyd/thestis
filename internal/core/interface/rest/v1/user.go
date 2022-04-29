package v1

import (
	"net/http"

	"github.com/harpyd/thestis/internal/core/interface/rest"
)

func authorize(w http.ResponseWriter, r *http.Request) (rest.User, bool) {
	user, err := rest.UserFromCtx(r.Context())
	if err != nil {
		rest.Unauthorized(string(ErrorSlugUnauthorizedUser), err, w, r)

		return rest.User{}, false
	}

	return user, true
}
