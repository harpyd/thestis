package v1

import (
	stdhttp "net/http"

	"github.com/harpyd/thestis/internal/port/http"
)

func unmarshalUser(w stdhttp.ResponseWriter, r *stdhttp.Request) (http.User, bool) {
	user, err := http.UserFromCtx(r.Context())
	if err != nil {
		http.Unauthorized(string(ErrorSlugUnauthorizedUser), err, w, r)

		return http.User{}, false
	}

	return user, true
}
