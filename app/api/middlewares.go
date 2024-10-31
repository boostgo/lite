package api

import (
	"context"
	"github.com/boostgo/lite/types/to"
	"github.com/labstack/echo/v4"
)

const (
	rawResponseKey = "lite-response-raw"
)

func Raw() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			localCtx := Context(ctx)
			localCtx = context.WithValue(localCtx, rawResponseKey, true)
			ctx.SetRequest(ctx.Request().WithContext(localCtx))
			return next(ctx)
		}
	}
}

func isRaw(ctx echo.Context) bool {
	return to.Bool(Context(ctx).Value(rawResponseKey))
}
