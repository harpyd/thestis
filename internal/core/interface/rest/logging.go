package rest

import (
	"bytes"
	"fmt"
	"io"
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
		service.StringLogField("requestId", requestID),
		service.StringLogField("httpMethod", r.Method),
		service.StringLogField("removeAddress", r.RemoteAddr),
		service.StringLogField("uri", r.RequestURI),
		service.StringLogField("requestBody", copyRequestBody(r)),
	)

	return &logEntry{
		logging:   l.logging,
		uri:       r.RequestURI,
		requestID: requestID,
	}
}

func copyRequestBody(r *http.Request) string {
	body, _ := io.ReadAll(r.Body)
	_ = r.Body.Close()
	r.Body = io.NopCloser(bytes.NewBuffer(body))

	return string(body)
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
		service.StringLogField("uri", e.uri),
		service.StringLogField("requestId", e.requestID),
		service.IntLogField("responseStatus", status),
		service.IntLogField("responseBytesLength", bytes),
		service.DurationLogField("responseElapsed", elapsed.Round(time.Millisecond/responseRounding)),
		service.DurationLogField("responseElapsed", elapsed.Round(time.Millisecond/responseRounding)),
	)
}

func (e *logEntry) Panic(v interface{}, stack []byte) {
	e.logging = e.logging.With(
		service.BytesLogField("stack", stack),
		service.StringLogField("panic", fmt.Sprintf("%+v", v)),
	)
}

func logger(r *http.Request) service.Logger {
	entry, ok := middleware.GetLogEntry(r).(*logEntry)
	if !ok {
		panic("LogEntry isn't *logEntry")
	}

	return entry.logging
}

func LoggingMiddleware(loggingService service.Logger) Middleware {
	return middleware.RequestLogger(&Formatter{logging: loggingService})
}
