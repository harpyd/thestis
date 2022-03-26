package cors

import (
	"net/http"

	"github.com/go-chi/cors"
)

const maxAge = 300

func Middleware(allowedOrigins []string) func(next http.Handler) http.Handler {
	return cors.New(cors.Options{
		AllowedOrigins:   allowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           maxAge,
	}).Handler
}
