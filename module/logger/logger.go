package logger

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/pohara/config"
	debugtrace "github.com/mrrizkin/pohara/module/debug-trace"
	"github.com/mrrizkin/pohara/module/logger/provider"
)

type LoggerProvider interface {
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
	Fatal(msg string, args ...interface{})
	FxLogEvent(event fxevent.Event)
	GetLogger() interface{}
}

type Logger struct {
	provider LoggerProvider
}

type Dependencies struct {
	fx.In

	Config *config.App
}

type Result struct {
	fx.Out

	Logger *Logger
}

func New(deps Dependencies) (Result, error) {
	p, err := provider.Zerolog(deps.Config)
	if err != nil {
		return Result{}, err
	}

	return Result{
		Logger: &Logger{
			provider: p,
		},
	}, nil
}

// usage
func (log *Logger) Info(msg string, args ...interface{}) {
	log.provider.Info(msg, args...)
}

func (log *Logger) Warn(msg string, args ...interface{}) {
	log.provider.Warn(msg, args...)
}

func (log *Logger) Error(msg string, args ...interface{}) {
	args = append(args, "stack_trace", debugtrace.StackTrace(9))
	log.provider.Error(msg, args...)
}

func (log *Logger) Fatal(msg string, args ...interface{}) {
	args = append(args, "stack_trace", debugtrace.StackTrace(9))
	log.provider.Fatal(msg, args...)
}

func (log *Logger) GetLogger() interface{} {
	return log.provider.GetLogger()
}

func (log *Logger) LogEvent(event fxevent.Event) {
	log.provider.FxLogEvent(event)
}
