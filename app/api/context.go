package api

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/validator"
	"github.com/boostgo/lite/types/param"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"sync"
	"time"
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

func Param(ctx echo.Context, paramName string) (param.Param, error) {
	value := ctx.Param(paramName)
	if value == "" {
		return param.Param{}, errs.
			New("Path param is empty").
			SetError(errs.ErrUnprocessableEntity).
			AddContext("param-name", paramName)
	}

	return param.New(value), nil
}

func QueryParam(ctx echo.Context, queryParamName string) param.Param {
	value := ctx.QueryParam(queryParamName)
	if value == "" {
		return param.Param{}
	}

	return param.New(value)
}

func Parse(ctx echo.Context, export any) error {
	if err := ctx.Bind(export); err != nil {
		return errs.
			New("Parse request body").
			SetError(err, errs.ErrUnprocessableEntity)
	}

	return _validator.Struct(export)
}

func Body(ctx echo.Context) (body []byte, err error) {
	if ctx.Request().Body == nil {
		return nil, nil
	}

	body, err = io.ReadAll(ctx.Request().Body)
	if err != nil {
		return nil, errs.
			New("Parse request body").
			SetType("API").
			SetError(err)
	}

	return body, nil
}

func Context(ctx echo.Context) context.Context {
	return ctx.Request().Context()
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

func Header(ctx echo.Context, key string) string {
	return ctx.Request().Header.Get(key)
}

func SetHeader(ctx echo.Context, key, value string) {
	ctx.Response().Header().Set(key, value)
}

func Cookie(ctx echo.Context, key string) string {
	cookie, err := ctx.Request().Cookie(key)
	if err != nil {
		return ""
	}

	return cookie.Value
}

func SetCookie(ctx echo.Context, key, value string) {
	cookie := &http.Cookie{}
	cookie.Name = key
	cookie.Value = value
	cookie.Path = "/"
	cookie.HttpOnly = true
	cookie.Secure = true
	cookie.Expires = time.Now().Add(time.Hour * 24 * 7)
	ctx.SetCookie(cookie)
}
