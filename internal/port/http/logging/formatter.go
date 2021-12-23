package logging

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"

	"github.com/harpyd/thestis/internal/app"
)

type Formatter struct {
	logging app.LoggingService
}

func (l *Formatter) NewLogEntry(r *http.Request) middleware.LogEntry {
	requestID := middleware.GetReqID(r.Context())

	l.logging.Info(
		"Request started",
		app.LogField{Key: "requestId", Value: requestID},
		app.LogField{Key: "httpMethod", Value: r.Method},
		app.LogField{Key: "remoteAddress", Value: r.RemoteAddr},
		app.LogField{Key: "uri", Value: r.RequestURI},
		app.LogField{Key: "requestBody", Value: copyRequestBody(r)},
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
	logging   app.LoggingService
	uri       string
	requestID string
}

const responseRounding = 100

func (e *logEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	e.logging.Info(
		"Request completed",
		app.LogField{Key: "uri", Value: e.uri},
		app.LogField{Key: "requestId", Value: e.requestID},
		app.LogField{Key: "responseStatus", Value: strconv.Itoa(status)},
		app.LogField{Key: "responseBytesLength", Value: strconv.Itoa(bytes)},
		app.LogField{Key: "responseElapsed", Value: elapsed.Round(time.Millisecond / responseRounding).String()},
	)
}

func (e *logEntry) Panic(v interface{}, stack []byte) {
	e.logging = e.logging.With(
		app.LogField{Key: "stack", Value: string(stack)},
		app.LogField{Key: "panic", Value: fmt.Sprintf("%+v", v)},
	)
}

func Logger(r *http.Request) app.LoggingService {
	entry, ok := middleware.GetLogEntry(r).(*logEntry)
	if !ok {
		panic("LogEntry isn't *logEntry")
	}

	return entry.logging
}
