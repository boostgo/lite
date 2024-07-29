package log

import (
	"context"
	"github.com/boostgo/lite/system/life"
	"github.com/rs/zerolog/log"
)

func Debug(ctx context.Context, namespace ...string) Event {
	return newEvent(ctx, log.Debug(), namespace...)
}

func Info(ctx context.Context, namespace ...string) Event {
	return newEvent(ctx, log.Info(), namespace...)
}

func Warn(ctx context.Context, namespace ...string) Event {
	return newEvent(ctx, log.Warn(), namespace...)
}

func Error(ctx context.Context, namespace ...string) Event {
	return newEvent(ctx, log.Error(), namespace...)
}

func Fatal(ctx context.Context, namespace ...string) Event {
	defer life.Cancel()
	return newEvent(ctx, log.Error(), namespace...).Bool("fatal", true)
}

func Namespace(namespace string) Logger {
	return &namespaced{namespace: namespace}
}

type Logger interface {
	Debug(ctx context.Context) Event
	Info(ctx context.Context) Event
	Warn(ctx context.Context) Event
	Error(ctx context.Context) Event
	Fatal(ctx context.Context) Event
}

type namespaced struct {
	namespace string
}

func (logger namespaced) Debug(ctx context.Context) Event {
	return Debug(ctx, logger.namespace)
}

func (logger namespaced) Info(ctx context.Context) Event {
	return Info(ctx, logger.namespace)
}

func (logger namespaced) Warn(ctx context.Context) Event {
	return Warn(ctx, logger.namespace)
}

func (logger namespaced) Error(ctx context.Context) Event {
	return Error(ctx, logger.namespace)
}

func (logger namespaced) Fatal(ctx context.Context) Event {
	return Fatal(ctx, logger.namespace)
}
