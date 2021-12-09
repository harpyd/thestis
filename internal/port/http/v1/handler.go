package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type handler struct{}

func NewHandler(r chi.Router) http.Handler {
	return HandlerFromMux(handler{}, r)
}
