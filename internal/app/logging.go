package app

import (
	"strconv"
	"time"
)

type LoggingService interface {
	With(fields ...LogField) LoggingService
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warn(msg string, err error, fields ...LogField)
	Error(msg string, err error, fields ...LogField)
	Fatal(msg string, err error, fields ...LogField)
}

type LogField struct {
	key   string
	value string
}

func StringLogField(key, value string) LogField {
	return LogField{
		key:   key,
		value: value,
	}
}

func IntLogField(key string, value int) LogField {
	return LogField{
		key:   key,
		value: strconv.Itoa(value),
	}
}

func DurationLogField(key string, value time.Duration) LogField {
	return LogField{
		key:   key,
		value: value.String(),
	}
}

func BytesLogField(key string, value []byte) LogField {
	return LogField{
		key:   key,
		value: string(value),
	}
}

func (f LogField) Key() string {
	return f.key
}

func (f LogField) Value() string {
	return f.value
}
