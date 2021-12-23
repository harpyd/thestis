package zap

import (
	"go.uber.org/zap"

	"github.com/harpyd/thestis/internal/app"
)

type LoggingService struct {
	logger *zap.Logger
}

func NewLoggingService(logger *zap.Logger) *LoggingService {
	return &LoggingService{
		logger: logger,
	}
}

func (l LoggingService) With(fields ...app.LogField) app.LoggingService {
	return &LoggingService{
		logger: l.logger.With(mapToZapFields(fields)...),
	}
}

func (l LoggingService) Debug(msg string, fields ...app.LogField) {
	l.logger.Debug(msg, mapToZapFields(fields)...)
}

func (l LoggingService) Info(msg string, fields ...app.LogField) {
	l.logger.Info(msg, mapToZapFields(fields)...)
}

func (l LoggingService) Warn(msg string, err error, fields ...app.LogField) {
	l.logger.Warn(msg, mapToZapFieldsWithErr(err, fields)...)
}

func (l LoggingService) Error(msg string, err error, fields ...app.LogField) {
	l.logger.Error(msg, mapToZapFieldsWithErr(err, fields)...)
}

func (l LoggingService) Fatal(msg string, err error, fields ...app.LogField) {
	l.logger.Fatal(msg, mapToZapFieldsWithErr(err, fields)...)
}

func mapToZapFields(fields []app.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.String(f.Key, f.Value))
	}

	return zapFields
}

func mapToZapFieldsWithErr(err error, fields []app.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields)+1)
	zapFields = append(zapFields, zap.Error(err))

	for _, f := range fields {
		zapFields = append(zapFields, zap.String(f.Key, f.Value))
	}

	return zapFields
}
