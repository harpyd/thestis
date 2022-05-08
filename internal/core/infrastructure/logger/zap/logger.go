package zap

import (
	"github.com/harpyd/thestis/internal/core/app/service"
	"go.uber.org/zap"
)

type Logger struct {
	logger *zap.Logger
}

func NewLogger(logger *zap.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l Logger) With(fields ...service.LogField) service.Logger {
	return &Logger{
		logger: l.logger.With(mapToZapFields(fields)...),
	}
}

func (l Logger) Debug(msg string, fields ...service.LogField) {
	l.logger.Debug(msg, mapToZapFields(fields)...)
}

func (l Logger) Info(msg string, fields ...service.LogField) {
	l.logger.Info(msg, mapToZapFields(fields)...)
}

func (l Logger) Warn(msg string, err error, fields ...service.LogField) {
	l.logger.Warn(msg, mapToZapFieldsWithErr(err, fields)...)
}

func (l Logger) Error(msg string, err error, fields ...service.LogField) {
	l.logger.Error(msg, mapToZapFieldsWithErr(err, fields)...)
}

func (l Logger) Fatal(msg string, err error, fields ...service.LogField) {
	l.logger.Fatal(msg, mapToZapFieldsWithErr(err, fields)...)
}

func mapToZapFields(fields []service.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.String(f.Key(), f.Value()))
	}

	return zapFields
}

func mapToZapFieldsWithErr(err error, fields []service.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)+1)
	zapFields = append(zapFields, zap.Error(err))

	for _, f := range fields {
		zapFields = append(zapFields, zap.String(f.Key(), f.Value()))
	}

	return zapFields
}
