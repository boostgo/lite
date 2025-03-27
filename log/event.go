package log

import (
	"context"
	"time"

	"github.com/boostgo/convert"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/trace"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// Event represents a log event. It is instanced by one of the level method of
// Logger and finalized by the Msg or Msgf method.
//
// Notice: this is wrap over zerolog.Event
type Event interface {
	Ctx(ctx context.Context) Event
	Send()
	Any(key string, value any) Event
	Arr(key string, args ...any) Event
	Err(err error) Event
	Errs(key string, errors []error) Event
	Msg(message string) Event
	Msgf(format string, args ...any) Event
	Str(key string, val string) Event
	UUID(key string, id uuid.UUID) Event
	Strs(key string, values []string) Event
	Int(key string, val int) Event
	Int8(key string, value int8) Event
	Int16(key string, value int16) Event
	Int32(key string, value int32) Event
	Int64(key string, value int64) Event
	Ints(key string, values []int) Event
	Ints8(key string, values []int8) Event
	Ints16(key string, values []int16) Event
	Ints32(key string, values []int32) Event
	Ints64(key string, values []int64) Event
	Uint(key string, val uint) Event
	Uint8(key string, value uint8) Event
	Uint16(key string, value uint16) Event
	Uint32(key string, value uint32) Event
	Uint64(key string, value uint64) Event
	Uints(key string, values []uint) Event
	Uints8(key string, values []uint8) Event
	Uints16(key string, values []uint16) Event
	Uints32(key string, values []uint32) Event
	Uints64(key string, values []uint64) Event
	Float32(key string, value float32) Event
	Floats32(key string, values []float32) Event
	Float64(key string, value float64) Event
	Floats64(key string, values []float64) Event
	Bool(key string, val bool) Event
	Time(key string, val time.Time) Event
	Times(key string, value []time.Time) Event
	Duration(key string, val time.Duration) Event
	Durations(key string, value []time.Duration) Event
	Obj(key string, obj any) Event
	Bytes(key string, bytes []byte) Event
	Type(key string, obj any) Event
	Namespace(namespace string) Event
}

type event struct {
	inner *zerolog.Event
}

func newEvent(inner *zerolog.Event, ctx ...context.Context) Event {
	e := &event{
		inner: inner,
	}

	if len(ctx) > 0 && ctx[0] != nil {
		e.Ctx(ctx[0])
	}

	return e
}

func (e *event) Send() {
	e.inner.Send()
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

func (e *event) Arr(key string, args ...any) Event {
	if args == nil {
		return e
	}

	stringArgs := make([]string, len(args))
	for i, arg := range args {
		stringArgs[i] = convert.String(arg)
	}
	e.inner.Strs(key, stringArgs)
	return e
}

func (e *event) Err(err error) Event {
	if err == nil {
		return e
	}

	custom, ok := errs.TryGet(err)
	if !ok {
		e.inner.Err(err)
	} else {
		if custom.Type() != "" {
			e.Str("error_type", custom.Type())
		}

		if custom.InnerError() != nil {
			e.Str("inner_error", custom.InnerError().Error())
		}

		if custom.Context() != nil && len(custom.Context()) > 0 {
			for key, value := range custom.Context() {
				e.Obj(key, value)
			}
		}

		e.Str("message", custom.Message())
	}
	return e
}

func (e *event) Errs(key string, errors []error) Event {
	e.inner.Errs(key, errors)
	return e
}

func (e *event) Msg(message string) Event {
	e.inner.Msg(message)
	return e
}

func (e *event) Msgf(format string, args ...any) Event {
	e.inner.Msgf(format, args...)
	return e
}

func (e *event) Str(key, value string) Event {
	e.inner.Str(key, value)
	return e
}

func (e *event) UUID(key string, id uuid.UUID) Event {
	e.inner.Str(key, id.String())
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

func (e *event) Int8(key string, value int8) Event {
	e.inner.Int8(key, value)
	return e
}

func (e *event) Int16(key string, value int16) Event {
	e.inner.Int16(key, value)
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

func (e *event) Ints8(key string, values []int8) Event {
	e.inner.Ints8(key, values)
	return e
}

func (e *event) Ints16(key string, values []int16) Event {
	e.inner.Ints16(key, values)
	return e
}

func (e *event) Ints32(key string, values []int32) Event {
	e.inner.Ints32(key, values)
	return e
}

func (e *event) Ints64(key string, values []int64) Event {
	e.inner.Ints64(key, values)
	return e
}

func (e *event) Uint(key string, value uint) Event {
	e.inner.Uint(key, value)
	return e
}

func (e *event) Uint8(key string, value uint8) Event {
	e.inner.Uint8(key, value)
	return e
}

func (e *event) Uint16(key string, value uint16) Event {
	e.inner.Uint16(key, value)
	return e
}

func (e *event) Uint32(key string, value uint32) Event {
	e.inner.Uint32(key, value)
	return e
}

func (e *event) Uint64(key string, value uint64) Event {
	e.inner.Uint64(key, value)
	return e
}

func (e *event) Uints(key string, value []uint) Event {
	e.inner.Uints(key, value)
	return e
}

func (e *event) Uints8(key string, value []uint8) Event {
	e.inner.Uints8(key, value)
	return e
}

func (e *event) Uints16(key string, value []uint16) Event {
	e.inner.Uints16(key, value)
	return e
}

func (e *event) Uints32(key string, value []uint32) Event {
	e.inner.Uints32(key, value)
	return e
}

func (e *event) Uints64(key string, value []uint64) Event {
	e.inner.Uints64(key, value)
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

func (e *event) Times(key string, value []time.Time) Event {
	e.inner.Times(key, value)
	return e
}

func (e *event) Duration(key string, value time.Duration) Event {
	e.inner.Dur(key, value)
	return e
}

func (e *event) Durations(key string, value []time.Duration) Event {
	e.inner.Durs(key, value)
	return e
}

func (e *event) Obj(key string, obj any) Event {
	e.Any(key, convert.String(obj))
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
