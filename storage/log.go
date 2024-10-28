package storage

import (
	"context"
	"github.com/boostgo/lite/types/to"
)

const noLogKey = "LITE_NO_LOG"

func NoLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, noLogKey, true)
}

func IsNoLog(ctx context.Context) bool {
	return to.Bool(ctx.Value(noLogKey))
}
