package trace

import (
	"context"
	"github.com/google/uuid"
)

const (
	key = "lite_trace-id"
)

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

func GetUUID(ctx context.Context) uuid.UUID {
	traceID := ctx.Value(key)
	if traceID == nil {
		return uuid.UUID{}
	}

	switch tid := traceID.(type) {
	case string:
		uuidVer, err := uuid.Parse(tid)
		if err != nil {
			return uuid.UUID{}
		}

		return uuidVer
	case uuid.UUID:
		return tid
	default:
		return uuid.UUID{}
	}
}
