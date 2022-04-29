package v1

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"

	"github.com/harpyd/thestis/internal/core/app"
	"github.com/harpyd/thestis/internal/core/interface/http"
)

type handler struct {
	app    *app.Application
	logger app.Logger
}

func NewHandler(
	application *app.Application,
	logger app.Logger,
	middlewares ...http.Middleware,
) stdhttp.Handler {
	r := chi.NewRouter()
	for _, m := range middlewares {
		r.Use(m)
	}

	return HandlerFromMux(handler{
		app:    application,
		logger: logger,
	}, r)
}
