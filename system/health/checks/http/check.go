package http

import (
	"context"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/system/health"
	"github.com/boostgo/lite/web"
	"net/http"
	"time"
)

func New(serviceName, url string, timeout ...time.Duration) health.Checker {
	return health.NewChecker(serviceName, func(ctx context.Context) (status health.Status, err error) {
		request := web.
			R(ctx).
			Header("Connection", "close")
		if len(timeout) > 0 {
			request.Timeout(timeout[0])
		}

		response, err := request.GET(url)
		if err != nil {
			return status, errs.
				New("Check failed on calling request").
				SetError(err).
				AddContext("url", url).
				AddContext("service_name", serviceName)
		}

		if response.StatusCode() >= http.StatusInternalServerError {
			return status, errs.
				New("Check failed on response status code").
				AddContext("url", url).
				AddContext("service_name", serviceName).
				AddContext("status_code", response.StatusCode())
		}

		return health.Status{
			Status: health.StatusHealthy,
		}, nil
	})
}
