package log

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/types/to"
	"github.com/rs/zerolog"
	"time"
)

type Event interface {
	Ctx(ctx context.Context) Event
	Any(key string, value any) Event
	Err(err error) Event
	Errs(key string, errors []error) Event
	Msg(message string) Event
	Msgf(format string, args ...any) Event
	Str(key string, val string) Event
	Strs(key string, values []string) Event
	Int(key string, val int) Event
	Int32(key string, value int32) Event
	Int64(key string, value int64) Event
	Ints(key string, values []int) Event
	Float32(key string, value float32) Event
	Floats32(key string, values []float32) Event
	Float64(key string, value float64) Event
	Floats64(key string, values []float64) Event
	Bool(key string, val bool) Event
	Time(key string, val time.Time) Event
	Duration(key string, val time.Duration) Event
	Obj(key string, obj any) Event
	Bytes(key string, bytes []byte) Event
	Type(key string, obj any) Event
	Namespace(namespace string) Event
}

type event struct {
	inner  *zerolog.Event
	called byte
}

func newEvent(inner *zerolog.Event) Event {
	return &event{
		inner: inner,
	}
}

func (e *event) Ctx(ctx context.Context) Event {
	if ctx == nil {
		return e
	}

	e.inner.Ctx(ctx)

	traceID := trace.Get(ctx)
	if traceID != "" {
		e.Str("trace_id", traceID)
	}

	return e
}

func (e *event) Any(key string, object any) Event {
	e.inner.Interface(key, object)
	return e
}

func (e *event) Err(err error) Event {
	custom, ok := errs.TryGet(err)
	if !ok {
		e.inner.Err(err)
	} else {
		e.Str("errorType", custom.Type())

		if custom.InnerError() != nil {
			e.Str("innerError", custom.InnerError().Error())
		}

		if custom.Context() != nil && len(custom.Context()) > 0 {
			for key, value := range custom.Context() {
				e.Obj(key, value)
			}
		}

		e.Msg(custom.Message())
	}
	return e
}

func (e *event) Errs(key string, errors []error) Event {
	e.inner.Errs(key, errors)
	return e
}

func (e *event) Msg(message string) Event {
	if e.called == 1 {
		return e
	}

	e.inner.Msg(message)
	e.called = 1
	return e
}

func (e *event) Msgf(format string, args ...any) Event {
	if e.called == 1 {
		return e
	}

	e.inner.Msgf(format, args...)
	e.called = 1
	return e
}

func (e *event) Str(key, value string) Event {
	e.inner.Str(key, value)
	return e
}

func (e *event) Strs(key string, values []string) Event {
	e.inner.Strs(key, values)
	return e
}

func (e *event) Int(key string, value int) Event {
	e.inner.Int(key, value)
	return e
}

func (e *event) Int32(key string, value int32) Event {
	e.inner.Int32(key, value)
	return e
}

func (e *event) Int64(key string, value int64) Event {
	e.inner.Int64(key, value)
	return e
}

func (e *event) Ints(key string, values []int) Event {
	e.inner.Ints(key, values)
	return e
}

func (e *event) Float32(key string, value float32) Event {
	e.inner.Float32(key, value)
	return e
}

func (e *event) Floats32(key string, values []float32) Event {
	e.inner.Floats32(key, values)
	return e
}

func (e *event) Float64(key string, value float64) Event {
	e.inner.Float64(key, value)
	return e
}

func (e *event) Floats64(key string, values []float64) Event {
	e.inner.Floats64(key, values)
	return e
}

func (e *event) Bool(key string, value bool) Event {
	e.inner.Bool(key, value)
	return e
}

func (e *event) Time(key string, value time.Time) Event {
	e.inner.Time(key, value)
	return e
}

func (e *event) Duration(key string, value time.Duration) Event {
	e.inner.Dur(key, value)
	return e
}

func (e *event) Obj(key string, obj any) Event {
	e.Any(key, to.String(obj))
	return e
}

func (e *event) Bytes(key string, bytes []byte) Event {
	e.inner.Bytes(key, bytes)
	return e
}

func (e *event) Type(key string, obj any) Event {
	e.inner.Type(key, obj)
	return e
}

func (e *event) Namespace(namespace string) Event {
	if namespace == "" {
		return e
	}

	e.Str("namespace", namespace)
	return e
}
