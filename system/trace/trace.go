package trace

import (
	"context"
	"github.com/google/uuid"
	"sync/atomic"
)

const (
	defaultKey = "X-Lite-Trace-ID"
)

var (
	_key = defaultKey
)

var (
	_masterMode = atomic.Bool{}
)

// IAmMaster set flag that current app is "Trace Master".
//
// "Trace Master" means that app will generate trace ids on every handler request, on kafka, rmq messages, etc...
func IAmMaster() {
	_masterMode.Store(true)
}

// AmIMaster returns flag if current app is "Trace Master".
//
// "Trace Master" means that app will generate trace ids on every handler request, on kafka, rmq messages, etc...
func AmIMaster() bool {
	return _masterMode.Load()
}

// SetMasterKey sets new master key for any kind of resource (HTTP, Kafka, RMQ, etc...)
func SetMasterKey(key string) {
	_key = key
}

// Key returns trace id key from HTTP request, kafka/rmq message, etc...
func Key() string {
	return _key
}

// Set sets trace id to new context
func Set(ctx context.Context, id string) context.Context {
	if Get(ctx) != "" {
		return ctx
	}

	return context.WithValue(ctx, _key, id)
}

// Get returns trace id from provided context
func Get(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	traceID := ctx.Value(_key)
	if traceID == nil {
		return ""
	}

	switch tid := traceID.(type) {
	case string:
		return tid
	case uuid.UUID:
		return tid.String()
	default:
		return ""
	}
}

// Has checks if provided context contain trace id
func Has(ctx context.Context) bool {
	return Get(ctx) != ""
}

// GetUUID returns trace id as "UUID" from context
func GetUUID(ctx context.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}

	traceID := ctx.Value(_key)
	if traceID == nil {
		return uuid.Nil
	}

	switch tid := traceID.(type) {
	case string:
		uuidVer, err := uuid.Parse(tid)
		if err != nil {
			return uuid.Nil
		}

		return uuidVer
	case uuid.UUID:
		return tid
	default:
		return uuid.Nil
	}
}

// ID generate new trace id as UUID
func ID() uuid.UUID {
	return uuid.New()
}

// String generate new trace id as string
func String() string {
	return ID().String()
}
