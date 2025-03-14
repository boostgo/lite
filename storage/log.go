package storage

import (
	"context"
	"github.com/boostgo/lite/types/to"
)

const noLogKey = "LITE_NO_LOG"

// NoLog set to context "no log" key.
//
// If key is set, logging queries will be turned off
func NoLog(ctx context.Context) context.Context {
	return context.WithValue(ctx, noLogKey, true)
}

// IsNoLog checks if context contain "no log" key
func IsNoLog(ctx context.Context) bool {
	return to.Bool(ctx.Value(noLogKey))
}
