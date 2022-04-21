package v1

import (
	stdhttp "net/http"

	"github.com/go-chi/render"

	"github.com/harpyd/thestis/internal/interface/http"
)

func decode(w stdhttp.ResponseWriter, r *stdhttp.Request, v interface{}) bool {
	if err := render.Decode(r, v); err != nil {
		http.BadRequest(string(ErrorSlugBadRequest), err, w, r)

		return false
	}

	return true
}
