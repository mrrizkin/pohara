package logger

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"

	"github.com/mrrizkin/pohara/config"
	"github.com/mrrizkin/pohara/internal/infrastructure/logger/util"
	"github.com/mrrizkin/pohara/internal/ports"
)

type ZeroLog struct {
	*zerolog.Logger
}

var Module = fx.Module("zerolog",
	fx.Provide(NewZeroLog),
	fx.Provide(func(logger *ZeroLog) ports.Logger { return logger }),
)

func NewZeroLog(config *config.App) (*ZeroLog, error) {
	var writers []io.Writer

	if config.LOG_CONSOLE {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	if config.LOG_FILE {
		rf, err := util.RollingFile(config)
		if err != nil {
			return nil, err
		}
		writers = append(writers, rf)
	}
	mw := io.MultiWriter(writers...)

	switch config.LOG_LEVEL {
	case "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "trace":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "disable":
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	zlogInstance := zerolog.New(mw).With().Timestamp().Logger()
	logger := ZeroLog{&zlogInstance}

	logger.Info(
		"logging configured",
		"fileLogging", config.LOG_FILE,
		"jsonLogOutput", config.LOG_JSON,
		"logDirectory", config.LOG_DIR,
		"fileName", config.APP_NAME+".log",
		"maxSizeMB", config.LOG_MAX_SIZE,
		"maxBackups", config.LOG_MAX_BACKUP,
		"maxAgeInDays", config.LOG_MAX_AGE,
	)

	return &logger, nil
}

// usage
func (z *ZeroLog) Scope(scope string) ports.Logger {
	zlog := z.Logger.With().Str("scope", scope).Logger()
	return &ZeroLog{
		Logger: &zlog,
	}
}

func (z *ZeroLog) Info(msg string, args ...interface{}) {
	z.argsParser(z.Logger.Info(), args...).Msg(msg)
}

func (z *ZeroLog) Warn(msg string, args ...interface{}) {
	z.argsParser(z.Logger.Warn(), args...).Msg(msg)
}

func (z *ZeroLog) Error(msg string, args ...interface{}) {
	z.argsParser(z.Logger.Error(), args...).Msg(msg)
}

func (z *ZeroLog) Fatal(msg string, args ...interface{}) {
	z.argsParser(z.Logger.Fatal(), args...).Msg(msg)
}

// LogEvent logs the given event to the provided Zerolog.
func (z *ZeroLog) LogEvent(event fxevent.Event) {
	switch e := event.(type) {
	case *fxevent.OnStartExecuting:
		z.Info(
			"OnStart hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)

	case *fxevent.OnStartExecuted:
		if e.Err != nil {
			z.Error(
				"OnStart hook failed",
				"error", e.Err,
				"callee", e.FunctionName,
				"caller", e.CallerName,
			)
		} else {
			z.Info(
				"OnStart hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.OnStopExecuting:
		z.Info(
			"OnStop hook executing",
			"callee", e.FunctionName,
			"caller", e.CallerName,
		)
	case *fxevent.OnStopExecuted:
		if e.Err != nil {
			z.Error(
				"OnStop hook failed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"error", e.Err,
			)
		} else {
			z.Info(
				"OnStop hook executed",
				"callee", e.FunctionName,
				"caller", e.CallerName,
				"runtime", e.Runtime.String(),
			)
		}
	case *fxevent.Supplied:
		if e.Err != nil {
			z.Error(
				"supplied",
				"error", e.Err,
				"type", e.TypeName,
				"module", e.ModuleName,
			)
		} else {
			z.Info(
				"supplied",
				"type", e.TypeName,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Provided:
		for _, rtype := range e.OutputTypeNames {
			z.Info(
				"provided",
				"constructor", e.ConstructorName,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			z.Error(
				"error encountered while applying options",
				"error", e.Err,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Decorated:
		for _, rtype := range e.OutputTypeNames {
			z.Info(
				"decorated",
				"decorator", e.DecoratorName,
				"module", e.ModuleName,
				"type", rtype,
			)
		}
		if e.Err != nil {
			z.Error(
				"error encountered while applying options",
				"error", e.Err,
				"module", e.ModuleName,
			)
		}
	case *fxevent.Invoking:
		// Do not log stack as it will make logs hard to read.
		z.Info(
			"invoking",
			"function", e.FunctionName,
			"module", e.ModuleName,
		)
	case *fxevent.Invoked:
		if e.Err != nil {
			z.Error(
				"invoke failed",
				"error", e.Err,
				"stack", e.Trace,
				"function", e.FunctionName,
			)
		}
	case *fxevent.Stopping:
		z.Info(
			"received signal",
			"signal", strings.ToUpper(e.Signal.String()),
		)
	case *fxevent.Stopped:
		if e.Err != nil {
			z.Error(
				"stop failed",
				"error", e.Err,
			)
		}
	case *fxevent.RollingBack:
		z.Error(
			"start failed, rolling back",
			"error", e.StartErr,
		)
	case *fxevent.RolledBack:
		if e.Err != nil {
			z.Error(
				"rollback failed",
				"error", e.Err,
			)
		}
	case *fxevent.Started:
		if e.Err != nil {
			z.Error(
				"start failed",
				"error", e.Err,
			)
		} else {
			z.Info("started")
		}
	case *fxevent.LoggerInitialized:
		if e.Err != nil {
			z.Error(
				"custom logger initialization failed",
				"error", e.Err,
			)
		} else {
			z.Info(
				"initialized custom fxevent.Logger",
				"function", e.ConstructorName,
			)
		}
	}
}

func (z *ZeroLog) GetLogger() interface{} {
	return z.Logger
}

func (z *ZeroLog) argsParser(event *zerolog.Event, args ...interface{}) *zerolog.Event {
	if len(args)%2 != 0 {
		z.Warn("logger: args don't match key val")
		return event
	}

	for i := 0; i < len(args); i += 2 {
		key, ok := args[i].(string)
		if !ok {
			z.Warn("info: non-string key provided")
			continue
		}

		switch value := args[i+1].(type) {
		case bool:
			event.Bool(key, value)
		case []bool:
			event.Bools(key, value)
		case string:
			event.Str(key, value)
		case []string:
			event.Strs(key, value)
		case int:
			event.Int(key, value)
		case []int:
			event.Ints(key, value)
		case int64:
			event.Int64(key, value)
		case []int64:
			event.Ints64(key, value)
		case float32:
			event.Float32(key, value)
		case float64:
			event.Float64(key, value)
		case time.Time:
			event.Time(key, value)
		case time.Duration:
			event.Dur(key, value)
		case []byte:
			event.Bytes(key, value)
		case error:
			event.Err(value)
		default:
			event.Interface(key, value)
		}
	}

	return event
}
