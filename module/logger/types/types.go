package types

import "go.uber.org/fx/fxevent"

type LoggerProvider interface {
	Scope(scope string) LoggerProvider
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	FxLogEvent(event fxevent.Event)
	GetLogger() interface{}
}
