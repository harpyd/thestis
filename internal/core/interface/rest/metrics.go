package rest

import (
	"net/http"
	"strconv"

	"github.com/urfave/negroni"
)

type MetricCollector interface {
	IncRequestsCount(status, method, path string)
}

func MetricsMiddleware(metric MetricCollector) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			method, path := r.Method, r.URL.Path

			lrw := negroni.NewResponseWriter(w)
			next.ServeHTTP(lrw, r)

			status := strconv.Itoa(lrw.Status())
			metric.IncRequestsCount(status, method, path)
		})
	}
}
