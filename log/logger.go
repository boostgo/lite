package log

import (
	"context"
	"github.com/boostgo/appx"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	_logger    = newLogger()
	_prettyLog = false
)

func newLogger() zerolog.Logger {
	switch os.Getenv("PRETTY_LOGGER") {
	case "true", "TRUE":
		_prettyLog = true
	}

	if _prettyLog {
		return log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	return zerolog.
		New(os.Stdout).
		With().
		Timestamp().
		Logger()
}

// PrettyLog enabled pretty logging mode.
//
// This mode could be activated by "PRETTY_LOGGER=true" env
func PrettyLog() {
	_prettyLog = true
	_logger = newLogger()
}

// Debug print log on debug level.
// Provided context use trace id
func Debug(ctx ...context.Context) Event {
	return newEvent(_logger.Debug(), ctx...)
}

// Info print log on info level.
// Provided context use trace id
func Info(ctx ...context.Context) Event {
	return newEvent(_logger.Info(), ctx...)
}

// Warn print log on warning level.
// Provided context use trace id
func Warn(ctx ...context.Context) Event {
	return newEvent(_logger.Warn(), ctx...)
}

// Error print log on error level.
// Provided context use trace id
func Error(ctx ...context.Context) Event {
	return newEvent(_logger.Error(), ctx...)
}

// Fatal print log on error level but with bool fatal=true.
// Provided context use trace id.
//
// Call AppCancel function
func Fatal(ctx ...context.Context) Event {
	defer appx.Cancel()
	return newEvent(_logger.Error().Bool("fatal", true), ctx...)
}

// Logger is wrap over zerolog logger
type Logger interface {
	// Debug print log on debug level.
	// Provided context use trace id
	Debug() Event
	// Info print log on info level.
	// Provided context use trace id
	Info() Event
	// Warn print log on warning level.
	// Provided context use trace id
	Warn() Event
	// Error print log on error level.
	// Provided context use trace id
	Error() Event
	// Fatal print log on error level but with bool fatal=true.
	// Provided context use trace id.
	//
	// Call life.Cancel() method which call graceful shutdown
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
