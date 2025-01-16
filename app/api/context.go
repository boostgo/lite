package api

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/validator"
	"github.com/boostgo/lite/types/format"
	"github.com/boostgo/lite/types/param"
	"github.com/labstack/echo/v4"
	"io"
	"net/http"
	"strings"
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

// Param returns [param.Param] object got from named path variable or not found param error.
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

// QueryParam returns query param variable as [param.Param] object or empty [param.Param] object if query param is not found.
func QueryParam(ctx echo.Context, queryParamName string) param.Param {
	value := ctx.QueryParam(queryParamName)
	if value == "" {
		return param.Empty()
	}

	return param.New(value)
}

// Parse try to parse request body to provided export object (must be pointer to structure object).
//
// After success parsing request body, run format converting (for "format" tags)
//
// After success format converting, run structure validation (for "validate" tags)
func Parse(ctx echo.Context, export any) error {
	if err := ctx.Bind(export); err != nil {
		return errs.
			New("Parse request body").
			SetError(err, errs.ErrUnprocessableEntity)
	}

	if err := format.Convert(export); err != nil {
		return err
	}

	return _validator.Struct(export)
}

// Body returns request body as []byte (slice of bytes)
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

// Context returns request context as context.Context object
func Context(ctx echo.Context) context.Context {
	return ctx.Request().Context()
}

// SetContext sets new context to echo.Context
func SetContext(ctx echo.Context, native context.Context) {
	ctx.SetRequest(ctx.Request().WithContext(native))
}

// File returns file as []byte (slice of bytes) from request by file name.
//
// Request body must be form data
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

// ParseForm get all form data object and convert them to map with [param.Param] objects.
//
// Notice: in this map no any files. Parse them by [File] function
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

// Header returns request header by provided key.
func Header(ctx echo.Context, key string) string {
	return ctx.Request().Header.Get(key)
}

// HeadersRaw return all headers as map with slice of values
func HeadersRaw(ctx echo.Context) map[string][]string {
	return ctx.Request().Header
}

// Headers return all headers as map with joined values
func Headers(ctx echo.Context) map[string]any {
	headers := make(map[string]any, len(ctx.Request().Header))
	for key, value := range ctx.Request().Header {
		headers[key] = strings.Join(value, ",")
	}
	return headers
}

// SetHeader sets new header to response
func SetHeader(ctx echo.Context, key, value string) {
	ctx.Response().Header().Set(key, value)
}

// Cookie returns request cookie by provided key
func Cookie(ctx echo.Context, key string) string {
	cookie, err := ctx.Request().Cookie(key)
	if err != nil {
		return ""
	}

	return cookie.Value
}

// CookiesRaw return all cookies as http.Cookie slice
func CookiesRaw(ctx echo.Context) []*http.Cookie {
	return ctx.Request().Cookies()
}

// Cookies return all cookies as map
func Cookies(ctx echo.Context) map[string]any {
	cookies := make(map[string]any)
	for _, cookie := range ctx.Request().Cookies() {
		cookies[cookie.Name] = cookie.Value
	}
	return cookies
}

// SetCookie sets new cookie to response
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
