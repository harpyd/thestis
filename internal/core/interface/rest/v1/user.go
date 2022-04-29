package v1

import (
	stdhttp "net/http"

	"github.com/harpyd/thestis/internal/core/interface/rest"
)

func authorize(w stdhttp.ResponseWriter, r *stdhttp.Request) (rest.User, bool) {
	user, err := rest.UserFromCtx(r.Context())
	if err != nil {
		rest.Unauthorized(string(ErrorSlugUnauthorizedUser), err, w, r)

		return rest.User{}, false
	}

	return user, true
}
