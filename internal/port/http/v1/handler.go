package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/harpyd/thestis/internal/app"
)

type handler struct {
	app *app.Application
}

func NewHandler(application *app.Application, r chi.Router) http.Handler {
	return HandlerFromMux(handler{
		app: application,
	}, r)
}
