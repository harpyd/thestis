package v1

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/harpyd/thestis/internal/core/adapter/driver/rest"
)

func decode(w http.ResponseWriter, r *http.Request, v interface{}) bool {
	if err := render.Decode(r, v); err != nil {
		rest.BadRequest(string(ErrorSlugBadRequest), err, w, r)

		return false
	}

	return true
}
