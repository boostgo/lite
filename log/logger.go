package log

import (
	"context"
	"github.com/boostgo/lite/config"
	"github.com/boostgo/lite/system/life"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

var (
	_logger    = newLogger()
	_prettyLog = false
)

func newLogger() zerolog.Logger {
	if config.Get("PRETTY_LOGGER").Bool() || _prettyLog {
		return log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger()
}

func PrettyLog() {
	_prettyLog = true
	_logger = newLogger()
}

// Debug print log on debug level.
// Provided context use trace id
func Debug(ctx ...context.Context) Event {
	return newEvent(_logger.Debug(), ctx...)
}

// Info print log on debug level.
// Provided context use trace id
func Info(ctx ...context.Context) Event {
	return newEvent(_logger.Info(), ctx...)
}

// Warn print log on debug level.
// Provided context use trace id
func Warn(ctx ...context.Context) Event {
	return newEvent(_logger.Warn(), ctx...)
}

// Error print log on debug level.
// Provided context use trace id
func Error(ctx ...context.Context) Event {
	return newEvent(_logger.Error(), ctx...)
}

// Fatal print log on debug level.
// Provided context use trace id.
// Call life.Cancel() method which call graceful shutdown
func Fatal(ctx ...context.Context) Event {
	defer life.Cancel()
	return newEvent(_logger.Error().Bool("fatal", true), ctx...)
}

// Logger wrap interface for zerolog logger
type Logger interface {
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
}

type wrapper struct {
	namespace string
	ctx       context.Context
}

// Namespace creates Logger implementation with namespace
func Namespace(namespace string) Logger {
	return Context(context.Background(), namespace)
}

// Context creates Logger implementation with context & namespace
func Context(ctx context.Context, namespace string) Logger {
	return &wrapper{
		ctx:       ctx,
		namespace: namespace,
	}
}

func (logger *wrapper) Debug() Event {
	return Debug().Ctx(logger.ctx).Namespace(logger.namespace)
}

func (logger *wrapper) Info() Event {
	return Info().Ctx(logger.ctx).Namespace(logger.namespace)
}

func (logger *wrapper) Warn() Event {
	return Warn().Ctx(logger.ctx).Namespace(logger.namespace)
}

func (logger *wrapper) Error() Event {
	return Error().Ctx(logger.ctx).Namespace(logger.namespace)
}

func (logger *wrapper) Fatal() Event {
	return Fatal().Ctx(logger.ctx).Namespace(logger.namespace)
}
