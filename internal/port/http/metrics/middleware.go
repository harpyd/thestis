package metrics

import (
	"net/http"
	"strconv"

	"github.com/urfave/negroni"

	"github.com/harpyd/thestis/internal/app"
)

func Middleware(metricsService app.MetricsService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method, path := r.Method, r.URL.Path

			lrw := negroni.NewResponseWriter(w)
			next.ServeHTTP(lrw, r)

			status := strconv.Itoa(lrw.Status())
			metricsService.IncRequestsCount(status, method, path)
		})
	}
}
