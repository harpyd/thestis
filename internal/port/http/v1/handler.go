package v1

import (
	stdhttp "net/http"

	"github.com/go-chi/chi/v5"

	"github.com/harpyd/thestis/internal/app"
	"github.com/harpyd/thestis/internal/port/http"
)

type handler struct {
	app *app.Application
}

func NewHandler(application *app.Application, middlewares ...http.Middleware) stdhttp.Handler {
	r := chi.NewRouter()
	for _, m := range middlewares {
		r.Use(m)
	}

	return HandlerFromMux(handler{
		app: application,
	}, r)
}
