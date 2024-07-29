package api

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/validator"
	"github.com/boostgo/lite/types/param"
	"github.com/labstack/echo/v4"
	"sync"
)

var (
	_validator     *validator.Validator
	_validatorOnce sync.Once
)

func init() {
	_validatorOnce.Do(func() {
		_validator, _ = validator.New()
	})
}

func Parse(ctx echo.Context, export any) error {
	if err := ctx.Bind(export); err != nil {
		return errs.
			New("Parse request body error").
			SetError(err, errs.ErrBadRequest)
	}

	contentType := ctx.Request().Header.Get("Content-Type")
	if contentType == "application/json" || contentType == "application/xml" {
		return _validator.Struct(export)
	}

	return nil
}

func Context(ctx echo.Context) context.Context {
	return ctx.Request().Context()
}

func QueryParam(ctx echo.Context, name string) param.Param {
	return param.New(ctx.QueryParam(name))
}
