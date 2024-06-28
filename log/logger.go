package log

import (
	"github.com/boostgo/lite/system/life"
	"github.com/rs/zerolog/log"
)

func Debug(namespace ...string) Event {
	return newEvent(log.Debug(), namespace...)
}

func Info(namespace ...string) Event {
	return newEvent(log.Info(), namespace...)
}

func Warn(namespace ...string) Event {
	return newEvent(log.Warn(), namespace...)
}

func Error(namespace ...string) Event {
	return newEvent(log.Error(), namespace...)
}

func Fatal(namespace ...string) Event {
	defer life.Cancel()
	return newEvent(log.Error(), namespace...).Bool("fatal", true)
}

func Namespace(namespace string) Logger {
	return &namespaced{namespace: namespace}
}

type Logger interface {
	Debug() Event
	Info() Event
	Warn() Event
	Error() Event
	Fatal() Event
}

type namespaced struct {
	namespace string
}

func (logger namespaced) Debug() Event {
	return Debug(logger.namespace)
}

func (logger namespaced) Info() Event {
	return Info(logger.namespace)
}

func (logger namespaced) Warn() Event {
	return Warn(logger.namespace)
}

func (logger namespaced) Error() Event {
	return Error(logger.namespace)
}

func (logger namespaced) Fatal() Event {
	return Fatal(logger.namespace)
}
