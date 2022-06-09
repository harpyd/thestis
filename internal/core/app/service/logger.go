package service

type Logger interface {
	With(args ...interface{}) Logger
	Named(name string) Logger
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
}
