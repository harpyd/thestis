package logging

import (
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/harpyd/thestis/internal/app"
)

func Middleware(loggingService app.LoggingService) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&Formatter{logging: loggingService})
}
