package trace

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

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

func FromResponse(response *http.Response) string {
	if response == nil {
		return ""
	}

	return response.Header.Get(key)
}

func SetRequest(request *http.Request, traceID string) {
	if request == nil {
		return
	}

	request.Header.Set(key, traceID)
}

func SetRequestCtx(ctx context.Context, request *http.Request) {
	traceID := Get(ctx)
	if traceID == "" {
		return
	}

	SetRequest(request, traceID)
}

func SetEchoCtx(ctx echo.Context, traceID string) {
	if ctx == nil {
		return
	}

	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), key, traceID)))
}

func GetEchoCtx(ctx echo.Context) string {
	return Get(ctx.Request().Context())
}
