package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type (
	Params struct {
		Middlewares []Middleware
		Routes      []Route
	}

	Route struct {
		Pattern string
		Handler http.Handler
	}

	Middleware func(next http.Handler) http.Handler
)

func NewHandler(params Params) http.Handler {
	rootRouter := chi.NewRouter()
	for _, m := range params.Middlewares {
		rootRouter.Use(m)
	}

	for _, route := range params.Routes {
		rootRouter.Mount(route.Pattern, route.Handler)
	}

	return rootRouter
}
