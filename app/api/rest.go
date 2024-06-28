package api

import (
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/types/flex"
	"github.com/boostgo/lite/types/to"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Ok(ctx echo.Context, body ...any) error {
	if len(body) == 0 {
		return ctx.String(http.StatusOK, "")
	}

	if flex.Type(body[0]).IsPrimitive() {
		return ctx.String(http.StatusOK, to.String(body[0]))
	}

	return ctx.JSON(http.StatusOK, body[0])
}

func Error(ctx echo.Context, err error) error {
	var body []byte
	status := http.StatusInternalServerError

	custom, ok := errs.TryGet(err)
	if ok {
		body = custom.JSON()
		status = custom.HttpCode()
	} else {
		body = errs.New(err.Error()).JSON()
	}

	return ctx.JSONBlob(status, body)
}

func Created(ctx echo.Context, body ...any) error {
	if len(body) == 0 {
		return ctx.String(http.StatusCreated, "")
	}

	switch value := body[0].(type) {
	case string, uuid.UUID: // provided id
		return ctx.JSON(http.StatusCreated, newCreatedID(value))
	default: // provided body
		return ctx.JSON(http.StatusCreated, body[0])
	}
}

func BadRequest(ctx echo.Context, message string) error {
	return Error(ctx, errs.New(message).SetHttpCode(http.StatusBadRequest))
}

func Unauthorized(ctx echo.Context, message ...string) error {
	return Error(ctx, errs.New("Unauthorized").SetHttpCode(http.StatusUnauthorized))
}

func Forbidden(ctx echo.Context) error {
	return Error(ctx, errs.New("Forbidden").SetHttpCode(http.StatusForbidden))
}

func NotFound(ctx echo.Context, message string) error {
	return Error(ctx, errs.New(message).SetHttpCode(http.StatusNotFound))
}

func Timeout(ctx echo.Context) error {
	return Error(ctx, errs.New("Request reached timeout").SetHttpCode(http.StatusRequestTimeout))
}

func MethodNotAllowed(ctx echo.Context) error {
	return Error(ctx, errs.New("Method not allowed").SetHttpCode(http.StatusMethodNotAllowed))
}

func UnprocessableEntity(ctx echo.Context, message string) error {
	return Error(ctx, errs.New(message).SetHttpCode(http.StatusUnprocessableEntity))
}

func TooManyRequests(ctx echo.Context) error {
	return Error(ctx, errs.New("Too many requests").SetHttpCode(http.StatusTooManyRequests))
}
