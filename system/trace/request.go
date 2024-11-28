package trace

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

// FromRequest returns trace id from HTTP request
func FromRequest(request *http.Request) string {
	if request == nil {
		return ""
	}

	traceID := request.Header.Get(key)
	if traceID == "" {
		c, err := request.Cookie(key)
		if err == nil {
			traceID = c.Value
		}
	}

	return traceID
}

// FromResponse returns trace id from HTTP response
func FromResponse(response *http.Response) string {
	if response == nil {
		return ""
	}

	return response.Header.Get(key)
}

// SetRequest sets trace id to HTTP request
func SetRequest(request *http.Request, traceID string) {
	if request == nil {
		return
	}

	request.Header.Set(key, traceID)
}

// SetRequestCtx sets trace id from context to HTTP request
func SetRequestCtx(ctx context.Context, request *http.Request) {
	traceID := Get(ctx)
	if traceID == "" {
		return
	}

	SetRequest(request, traceID)
}

// SetEchoCtx sets trace id to HTTP request into echo context
func SetEchoCtx(ctx echo.Context, traceID string) {
	if ctx == nil {
		return
	}

	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), key, traceID)))
}

// GetEchoCtx returns trace id from request in echo context
func GetEchoCtx(ctx echo.Context) string {
	return Get(ctx.Request().Context())
}
