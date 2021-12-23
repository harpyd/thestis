package app

type LoggingService interface {
	With(fields ...LogField) LoggingService
	Debug(msg string, fields ...LogField)
	Info(msg string, fields ...LogField)
	Warn(msg string, err error, fields ...LogField)
	Error(msg string, err error, fields ...LogField)
	Fatal(msg string, err error, fields ...LogField)
}

type LogField struct {
	Key   string
	Value string
}
