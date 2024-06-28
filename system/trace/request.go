package trace

import (
	"context"
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
