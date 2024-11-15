package trace

import (
	"context"
	"github.com/google/uuid"
	"sync/atomic"
)

const (
	key = "X-Lite-Trace-ID"
)

var (
	_masterMode = atomic.Bool{}
)

func IAmMaster() {
	_masterMode.Store(true)
}

func AmIMaster() bool {
	return _masterMode.Load()
}

func Key() string {
	return key
}

func Set(ctx context.Context, id string) context.Context {
	if Get(ctx) != "" {
		return ctx
	}

	return context.WithValue(ctx, key, id)
}

func Get(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	traceID := ctx.Value(key)
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

func Has(ctx context.Context) bool {
	return Get(ctx) != ""
}

func GetUUID(ctx context.Context) uuid.UUID {
	if ctx == nil {
		return uuid.Nil
	}

	traceID := ctx.Value(key)
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

func ID() uuid.UUID {
	return uuid.New()
}

func String() string {
	return ID().String()
}
