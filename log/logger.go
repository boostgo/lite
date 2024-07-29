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

func Namespace(namespace string) Logger {
	return &wrapper{
		ctx:       context.Background(),
		namespace: namespace,
	}
}

func Context(ctx context.Context, namespace string) Logger {
	return &wrapper{
		ctx:       ctx,
		namespace: namespace,
	}
}

type Logger interface {
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
}

func New(ctx context.Context, namespace ...string) Logger {
	if len(namespace) > 0 && namespace[0] != "" {
		return Context(ctx, namespace[0])
	}

	return Context(ctx, "")
}

type wrapper struct {
	namespace string
	ctx       context.Context
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
