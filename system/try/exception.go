package try

import (
	"context"
	"errors"
	"runtime/debug"

	"github.com/boostgo/convert"
	"github.com/boostgo/lite/errs"
)

// Try recovers if panic was thrown.
//
// Return error of provided function and recover error
func Try(tryFunc func() error) (err error) {
	defer func() {
		if err == nil {
			err = CatchPanic(recover())
		}
	}()

	return tryFunc()
}

// Ctx is like Try but provided function has context as an argument
func Ctx(ctx context.Context, tryFunc func(ctx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return Try(func() error {
		return tryFunc(ctx)
	})
}

// Must run provided function but ignore error
func Must(tryFunc func() error) {
	_ = Try(tryFunc)
}

// CatchPanic got recover() return value and convert it to error
func CatchPanic(err any) error {
	if err == nil {
		return nil
	}

	return errs.
		New("PANIC RECOVER").
		SetError(errors.New(convert.String(err))).
		AddContext("trace", convert.String(debug.Stack()))
}
