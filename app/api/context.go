package api

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/validator"
	"github.com/boostgo/lite/types/param"
	"github.com/labstack/echo/v4"
	"io"
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

func Param(ctx echo.Context, name string) param.Param {
	return param.New(ctx.Param(name))
}

func File(ctx echo.Context, name string) (content []byte, err error) {
	defer errs.Wrap("API", &err, "Read form file error")

	header, err := ctx.FormFile(name)
	if err != nil {
		return content, err
	}

	file, err := header.Open()
	if err != nil {
		return content, err
	}
	defer file.Close()

	return io.ReadAll(file)
}

func ParseForm(ctx echo.Context) (map[string]param.Param, error) {
	form, err := ctx.MultipartForm()
	if err != nil {
		return nil, err
	}

	exportMap := make(map[string]param.Param)
	for key, values := range form.Value {
		if len(values) == 0 {
			continue
		}

		exportMap[key] = param.New(values[0])
	}

	return exportMap, nil
}

func Get[T any](ctx echo.Context, key string) T {
	return ctx.Get(key)
}
