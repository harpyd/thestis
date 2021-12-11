package http

import (
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"go.uber.org/zap"

	v1 "github.com/harpyd/thestis/internal/port/http/v1"
	"github.com/harpyd/thestis/pkg/logging"
)

func NewHandler(logger *zap.Logger) http.Handler {
	apiRouter := chi.NewRouter()
	addMiddlewares(apiRouter, logger)

	rootRouter := chi.NewRouter()
	rootRouter.Mount("/v1", v1.NewHandler(apiRouter))

	return rootRouter
}

func addMiddlewares(router *chi.Mux, logger *zap.Logger) {
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(logging.NewStructuredLogger(logger))
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
