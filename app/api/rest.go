package api

import (
	"encoding/json"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/types/flex"
	"github.com/boostgo/lite/types/to"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Failure(ctx echo.Context, status int, err error) error {
	const defaultErrorType = "ERROR"

	var output errorOutput
	output.Status = statusFailure

	custom, ok := errs.TryGet(err)
	if ok {
		output.Message = custom.Message()
		output.Type = custom.Type()
		output.Context = custom.Context()
	} else {
		output.Message = err.Error()
		output.Type = defaultErrorType
	}

	outputBlob, _ := json.Marshal(output)
	return ctx.JSONBlob(status, outputBlob)
}

func Error(ctx echo.Context, err error) error {
	return Failure(ctx, errStatusCode(err), err)
}

func Success(ctx echo.Context, status int, body ...any) error {
	if len(body) == 0 {
		return ctx.String(status, "")
	}

	if flex.Type(body[0]).IsPrimitive() {
		return ctx.String(status, to.String(body[0]))
	}

	return ctx.JSON(status, newSuccess(body[0]))
}

func Ok(ctx echo.Context, body ...any) error {
	return Success(ctx, http.StatusOK, body...)
}

func Created(ctx echo.Context, body ...any) error {
	switch value := body[0].(type) {
	case string, uuid.UUID: // provided id
		return Success(ctx, http.StatusCreated, newCreatedID(value))
	default: // provided body
		return Success(ctx, http.StatusCreated, newCreatedID(value))
	}
}
