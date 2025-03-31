package api

import (
	"encoding/json"
	"net/http"
	"reflect"

	"github.com/boostgo/convert"
	"github.com/boostgo/httpx"
	"github.com/boostgo/trace"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	traceKeySession = "bgo_trace_id"
	traceKeyHeader  = "X-Trace-ID"
	traceKeyCookie  = "trace_id"
)

// Failure returns response with some error status and convert provided error to
// errorOutput object and then convert it to JSON response.
//
// Sets trace id to the response if it was in request context.
//
// If error is custom from package "errs", output will build from custom error.
//
// If errors is custom and there is "trace" key in context, it will be ignored for outputError
func Failure(ctx echo.Context, status int, err error) error {
	// set trace ID
	traceID := trace.Get(Context(ctx), traceKeySession)
	if traceID != "" {
		ctx.Response().Header().Set(traceKeyHeader, traceID)
		ctx.Response().Header().Set("X-Request-ID", traceID)
	}

	response := httpx.NewFailureResponse(err)
	blob, _ := json.Marshal(response)

	// convert output object to bytes
	return ctx.JSONBlob(status, blob)
}

// Error is wrap function above [Failure] function with auto defining status code by provided error.
//
// There is a list of errors in "errs" packages and if provided error is one of them, it has own code representation
func Error(ctx echo.Context, err error) error {
	return Failure(ctx, httpx.StatusCodeByError(err), err)
}

// Success returns response with success code & successOutput object and convert it to JSON response.
//
// Sets trace id to the response if it was in request context.
//
// If provided body exist, and it is "primitive" response will be in raw (no successOutput object).
//
// If context contain "raw" middleware key, response will be in raw (no successOutput object).
//
// If body is not provided, will be returned empty string
func Success(ctx echo.Context, status int, body ...any) error {
	// set trace ID
	traceID := trace.Get(Context(ctx), traceKeySession)
	if traceID != "" {
		ctx.Response().Header().Set(traceKeyHeader, traceID)
		ctx.Response().Header().Set("X-Request-ID", traceID)
	}

	// return empty response if no response body
	if len(body) == 0 {
		return ctx.String(status, "")
	}

	if isPrimitive(body[0]) {
		return ctx.String(status, convert.String(body[0]))
	}

	if isRaw(ctx) {
		return ctx.JSON(status, body[0])
	}

	return ctx.JSON(status, httpx.NewSuccessResponse(body[0]))
}

// SuccessRaw returns response in "raw" way
func SuccessRaw(ctx echo.Context, status int, body []byte, contentType ...string) error {
	cType := httpx.ContentTypeBytes
	if len(contentType) > 0 && contentType[0] != "" {
		cType = contentType[0]
	}

	return ctx.Blob(status, cType, body)
}

// ReturnExcel returns response with Excel file content type
func ReturnExcel(ctx echo.Context, name string, file []byte) error {
	ctx.Response().Header().Set("Content-Disposition", "attachment; filename="+name)
	return ctx.Blob(http.StatusOK, httpx.ContentTypeExcel, file)
}

// Ok is wrap function over [Success] function.
//
// Sets HTTP code "OK" 200
func Ok(ctx echo.Context, body ...any) error {
	return Success(ctx, http.StatusOK, body...)
}

// OkRaw is wrap function over [SuccessRaw] function.
//
// Sets HTTP code "OK" 200
func OkRaw(ctx echo.Context, body []byte) error {
	return SuccessRaw(ctx, http.StatusOK, body)
}

// Created is wrap function over [Success] function.
//
// Sets HTTP code "Created" 201
func Created(ctx echo.Context, body ...any) error {
	if len(body) == 0 {
		return Success(ctx, http.StatusCreated)
	}

	switch value := body[0].(type) {
	case string, uuid.UUID, int, int64, int32: // provided id
		return Success(ctx, http.StatusCreated, httpx.NewCreatedResponse(value))
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
