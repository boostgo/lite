package trace

import (
	"context"
	"github.com/labstack/echo/v4"
	"net/http"
)

func FromRequest(request *http.Request) string {
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
	return response.Header.Get(key)
}

func SetRequest(request *http.Request, traceID string) {
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
	ctx.SetRequest(ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), key, traceID)))
}
