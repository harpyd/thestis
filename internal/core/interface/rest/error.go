package rest

import (
	"net/http"

	"github.com/go-chi/render"
)

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, http.StatusBadRequest)
}

func Unauthorized(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, http.StatusUnauthorized)
}

func Forbidden(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, http.StatusForbidden)
}

func NotFound(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, http.StatusNotFound)
}

func Conflict(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, http.StatusConflict)
}

func UnprocessableEntity(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, http.StatusUnprocessableEntity)
}

func InternalServerError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, http.StatusInternalServerError)
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, status int) {
	var details string
	if err != nil {
		details = err.Error()
	}

	resp := ErrorResponse{slug, details, status}

	if err := render.Render(w, r, resp); err != nil {
		panic(err)
	}
}

type ErrorResponse struct {
	Slug       string `json:"slug"`
	Details    string `json:"details"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, _ *http.Request) error {
	w.WriteHeader(e.httpStatus)

	return nil
}
