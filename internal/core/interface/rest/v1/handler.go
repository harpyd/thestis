package v1

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/internal/core/interface/rest"
)

type handler struct {
	app    *app.Application
	logger service.Logger
}

func NewHandler(
	application *app.Application,
	logger service.Logger,
	middlewares ...rest.Middleware,
) http.Handler {
	r := chi.NewRouter()
	for _, m := range middlewares {
		r.Use(m)
	}

	return HandlerFromMux(handler{
		app:    application,
		logger: logger,
	}, r)
}
