package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	"github.com/harpyd/thestis/internal/port/http/logging"
)

type Route struct {
	Pattern string
	Handler http.Handler
}

func NewHandler(logger *zap.Logger, routes ...Route) http.Handler {
	rootRouter := chi.NewRouter()
	addMiddlewares(rootRouter, logger)

	for _, route := range routes {
		rootRouter.Mount(route.Pattern, route.Handler)
	}

	return rootRouter
}

func addMiddlewares(router *chi.Mux, logger *zap.Logger) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logging.NewFormatter(logger))
	router.Use(middleware.Recoverer)
	addCORSMiddleware(router)
	router.Use(middleware.NoCache)
}

const maxAge = 300

func addCORSMiddleware(router *chi.Mux) {
	allowedOrigins := strings.Split(os.Getenv("CORS_ALLOWED_ORIGINS"), ";")
	if len(allowedOrigins) == 0 {
		return
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           maxAge,
	})
	router.Use(corsMiddleware.Handler)
}
