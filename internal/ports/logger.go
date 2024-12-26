package ports

import "go.uber.org/fx/fxevent"

type Logger interface {
	Scope(scope string) Logger
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	LogEvent(event fxevent.Event)
	GetLogger() interface{}
}
