package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/harpyd/thestis/internal/core/app/service"
)

type Formatter struct {
	logging service.Logger
}

func (l *Formatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	requestID := middleware.GetReqID(r.Context())

	l.logging.Info(
		"Request started",
		"requestId", requestID,
		"httpMethod", r.Method,
		"remoteAddress", r.RemoteAddr,
		"uri", r.RequestURI,
	)

	return &logEntry{
		logging:   l.logging,
		uri:       r.RequestURI,
		requestID: requestID,
	}
}

type logEntry struct {
	logging   service.Logger
	uri       string
	requestID string
}

const responseRounding = 100

func (e *logEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	e.logging.Info(
		"Request completed",
		"uri", e.uri,
		"requestId", e.requestID,
		"responseStatus", status,
		"responseBytesLength", bytes,
		"responseElapsed", elapsed.Round(time.Millisecond/responseRounding),
	)
}

func (e *logEntry) Panic(v interface{}, stack []byte) {
	e.logging = e.logging.With(
		"stack", stack,
		"panic", fmt.Sprintf("%+v", v),
	)
}

func LoggingMiddleware(logger service.Logger) Middleware {
	return middleware.RequestLogger(&Formatter{logging: logger.Named("API")})
}
