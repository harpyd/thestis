package rest

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/harpyd/thestis/internal/core/app"
)

type Formatter struct {
	logging app.Logger
}

func (l *Formatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	requestID := middleware.GetReqID(r.Context())

	l.logging.Info(
		"Request started",
		app.StringLogField("requestId", requestID),
		app.StringLogField("httpMethod", r.Method),
		app.StringLogField("removeAddress", r.RemoteAddr),
		app.StringLogField("uri", r.RequestURI),
		app.StringLogField("requestBody", copyRequestBody(r)),
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
	logging   app.Logger
	uri       string
	requestID string
}

const responseRounding = 100

func (e *logEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	e.logging.Info(
		"Request completed",
		app.StringLogField("uri", e.uri),
		app.StringLogField("requestId", e.requestID),
		app.IntLogField("responseStatus", status),
		app.IntLogField("responseBytesLength", bytes),
		app.DurationLogField("responseElapsed", elapsed.Round(time.Millisecond/responseRounding)),
		app.DurationLogField("responseElapsed", elapsed.Round(time.Millisecond/responseRounding)),
	)
}

func (e *logEntry) Panic(v interface{}, stack []byte) {
	e.logging = e.logging.With(
		app.BytesLogField("stack", stack),
		app.StringLogField("panic", fmt.Sprintf("%+v", v)),
	)
}

func logger(r *http.Request) app.Logger {
	entry, ok := middleware.GetLogEntry(r).(*logEntry)
	if !ok {
		panic("LogEntry isn't *logEntry")
	}

	return entry.logging
}

func LoggingMiddleware(loggingService app.Logger) Middleware {
	return middleware.RequestLogger(&Formatter{logging: loggingService})
}
