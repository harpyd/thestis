package zap

import (
	"go.uber.org/zap"

	"github.com/harpyd/thestis/internal/core/app/service"
)

type Logger struct {
	base *zap.SugaredLogger
}

func NewLogger(logger *zap.Logger) Logger {
	return Logger{
		base: logger.WithOptions(zap.AddCallerSkip(1)).Sugar(),
	}
}

func (l Logger) With(args ...interface{}) service.Logger {
	return Logger{
		base: l.base.With(args),
	}
}

func (l Logger) Named(name string) service.Logger {
	return Logger{
		base: l.base.Named(name),
	}
}

func (l Logger) Debug(msg string, args ...interface{}) {
	l.base.Debugw(msg, args...)
}

func (l Logger) Info(msg string, args ...interface{}) {
	l.base.Infow(msg, args...)
}

func (l Logger) Warn(msg string, args ...interface{}) {
	l.base.Warnw(msg, args...)
}

func (l Logger) Error(msg string, args ...interface{}) {
	l.base.Errorw(msg, args...)
}

func (l Logger) Fatal(msg string, args ...interface{}) {
	l.base.Fatalw(msg, args...)
}
