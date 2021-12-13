package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/harpyd/thestis/internal/app"
)

type handler struct {
	application app.Application
}

func NewHandler(application app.Application, r chi.Router) http.Handler {
	return HandlerFromMux(handler{
		application: application,
	}, r)
}
