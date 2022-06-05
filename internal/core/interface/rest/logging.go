package rest

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/harpyd/thestis/internal/core/app/service"
	"github.com/harpyd/thestis/pkg/correlationid"
)

type Formatter struct {
	logging service.Logger
}

func (l *Formatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	correlationID, _ := correlationid.FromCtx(r.Context())

	l.logging.Info(
		"Request started",
		"correlationId", correlationID,
		"httpMethod", r.Method,
		"remoteAddress", r.RemoteAddr,
		"uri", r.RequestURI,
	)

	return &logEntry{
		logging:       l.logging,
		uri:           r.RequestURI,
		correlationID: correlationID,
	}
}

type logEntry struct {
	logging       service.Logger
	uri           string
	correlationID string
}

const responseRounding = 100

func (e *logEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	e.logging.Info(
		"Request completed",
		"uri", e.uri,
		"correlationId", e.correlationID,
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
