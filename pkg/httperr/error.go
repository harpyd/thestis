package httperr

import (
	"net/http"

	"github.com/go-chi/render"
	"go.uber.org/zap"

	"github.com/harpyd/thestis/pkg/logging"
)

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Bad Request", http.StatusBadRequest)
}

func Unauthorized(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Unauthorised", http.StatusUnauthorized)
}

func Forbidden(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Forbidden", http.StatusForbidden)
}

func NotFound(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Not Found", http.StatusNotFound)
}

func UnprocessableEntity(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Unprocessable Entity", http.StatusUnprocessableEntity)
}

func InternalServerError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, "Internal Server Error", http.StatusInternalServerError)
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, logMSg string, status int) {
	logging.GetLogger(r).Warn(logMSg, zap.Error(err), zap.String("errorSlug", slug))

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
