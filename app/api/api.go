package api

import (
	"encoding/json"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/types/to"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
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
		if custom.InnerError() != nil {
			output.Inner = custom.InnerError().Error()
		}
	} else {
		output.Message = err.Error()
		output.Type = defaultErrorType
	}

	log.
		Context(Context(ctx), "API").
		Error().
		Int("status", status).
		Err(err).
		Msg("Failure request")

	outputBlob, _ := json.Marshal(output)
	return ctx.JSONBlob(status, outputBlob)
}

func Error(ctx echo.Context, err error) error {
	return Failure(ctx, errStatusCode(err), err)
}

func Success(ctx echo.Context, status int, body ...any) error {
	// set trace ID
	traceID := trace.Get(Context(ctx))
	if traceID != "" {
		ctx.Response().Header().Set(trace.Key(), traceID)
		ctx.Response().Header().Set("X-Request-ID", traceID)
	}

	// return empty response if no response body
	if len(body) == 0 {
		return ctx.String(status, "")
	}

	if isPrimitive(body[0]) {
		return ctx.String(status, to.String(body[0]))
	}

	if isRaw(ctx) {
		return ctx.JSON(status, body[0])
	}

	return ctx.JSON(status, newSuccess(body[0]))
}

func ReturnExcel(ctx echo.Context, name string, file []byte) error {
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename="+name)
	return ctx.Blob(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", file)
}

func Ok(ctx echo.Context, body ...any) error {
	return Success(ctx, http.StatusOK, body...)
}

func Created(ctx echo.Context, body ...any) error {
	if len(body) == 0 {
		return Success(ctx, http.StatusCreated)
	}

	switch value := body[0].(type) {
	case string, uuid.UUID, int, int64, int32: // provided id
		return Success(ctx, http.StatusCreated, newCreatedID(value))
	default: // provided body
		return Success(ctx, http.StatusCreated, value)
	}
}

func isPrimitive(object any) bool {
	switch reflect.TypeOf(object).Kind() {
	case reflect.Ptr, reflect.Struct, reflect.Interface,
		reflect.Slice, reflect.Array, reflect.Map:
		return false
	default:
		return true
	}
}
