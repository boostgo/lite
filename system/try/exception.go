package try

import (
	"context"
	"errors"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/types/to"
	"runtime/debug"
)

func Try(tryFunc func() error) (err error) {
	defer func() {
		if err == nil {
			err = CatchPanic(recover())
		}
	}()

	return tryFunc()
}

func Ctx(ctx context.Context, tryFunc func(ctx context.Context) error) error {
	if ctx == nil {
		ctx = context.Background()
	}

	return Try(func() error {
		return tryFunc(ctx)
	})
}

func Must(tryFunc func() error) {
	_ = Try(tryFunc)
}

func CatchPanic(err any) error {
	if err == nil {
		return nil
	}

	return errs.New("PANIC RECOVER").
		SetError(errors.New(to.String(err))).
		AddContext("trace", to.String(debug.Stack()))
}
