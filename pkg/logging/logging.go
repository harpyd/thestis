package logging

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type StructuredLogger struct {
	Logger *zap.Logger
}

func NewStructuredLogger(logger *zap.Logger) func(next http.Handler) http.Handler {
	return middleware.RequestLogger(&StructuredLogger{Logger: logger})
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	requestID := middleware.GetReqID(r.Context())

	l.Logger.Info(
		"Request started",
		zap.String("requestId", requestID),
		zap.String("httpMethod", r.Method),
		zap.String("remoteAddress", r.RemoteAddr),
		zap.String("uri", r.RequestURI),
		zap.String("requestBody", copyRequestBody(r)),
	)

	return &StructuredLoggerEntry{
		Logger:    l.Logger,
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

type StructuredLoggerEntry struct {
	Logger    *zap.Logger
	uri       string
	requestID string
}

const responseRounding = 100

func (e *StructuredLoggerEntry) Write(status, bytes int, _ http.Header, elapsed time.Duration, _ interface{}) {
	e.Logger.Info(
		"Request completed",
		zap.String("uri", e.uri),
		zap.String("requestId", e.requestID),
		zap.Int("responseStatus", status),
		zap.Int("responseBytesLength", bytes),
		zap.Duration("responseElapsed", elapsed.Round(time.Millisecond/responseRounding)),
	)
}

func (e *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	e.Logger = e.Logger.With(
		zap.String("stack", string(stack)),
		zap.String("panic", fmt.Sprintf("%+v", v)),
	)
}
