package log

import (
	"context"
	"github.com/boostgo/lite/system/life"
	"github.com/rs/zerolog/log"
)

func Debug() Event {
	return newEvent(log.Debug())
}

func Info() Event {
	return newEvent(log.Info())
}

func Warn() Event {
	return newEvent(log.Warn())
}

func Error() Event {
	return newEvent(log.Error())
}

func Fatal() Event {
	defer life.Cancel()
	return newEvent(log.Error().Bool("fatal", true))
}

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

func Namespace(namespace string) Logger {
	return Context(context.Background(), namespace)
}

func Context(ctx context.Context, namespace string) Logger {
	return &wrapper{
		ctx:       ctx,
		namespace: namespace,
	}
}

func (logger wrapper) Debug() Event {
	return Debug().Ctx(logger.ctx).Namespace(logger.namespace)
}

func (logger wrapper) Info() Event {
	return Info().Ctx(logger.ctx).Namespace(logger.namespace)
}

func (logger wrapper) Warn() Event {
	return Warn().Ctx(logger.ctx).Namespace(logger.namespace)
}

func (logger wrapper) Error() Event {
	return Error().Ctx(logger.ctx).Namespace(logger.namespace)
}

func (logger wrapper) Fatal() Event {
	return Fatal().Ctx(logger.ctx).Namespace(logger.namespace)
}
