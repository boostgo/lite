package gateway

import (
	"github.com/boostgo/lite/app/api"
	"github.com/boostgo/lite/errs"
	"github.com/labstack/echo/v4"
)

// Gateway controls all services and their routes
type Gateway struct {
	services []Service
}

// New creates [Gateway] object.
func New() *Gateway {
	return &Gateway{
		services: make([]Service, 0),
	}
}

// RegisterService append new services
func (gw *Gateway) RegisterService(services ...Service) *Gateway {
	if len(services) == 0 {
		return gw
	}

	gw.services = append(gw.services, services...)
	return gw
}

// Match finds matching [Route] and [Service] by request method & path.
func (gw *Gateway) Match(method, path string) (Service, Route, bool) {
	for _, s := range gw.services {
		r, match := s.Match(method, path)
		if !match {
			continue
		}

		return s, r, true
	}

	return nil, nil, false
}

// Handler generate [echo.Echo] handler for redirecting/proxying to match request
func (gw *Gateway) Handler() echo.HandlerFunc {
	return func(ctx echo.Context) error {
		matchService, matchRoute, match := gw.Match(ctx.Request().Method, ctx.Request().URL.Path)
		if !match {
			return api.Error(ctx, errs.
				New("No matching gateway route").
				SetError(errs.ErrNotFound).
				AddContext("method", ctx.Request().Method).
				AddContext("url", ctx.Request().URL.String()))
		}

		requestBody, err := api.Body(ctx)
		if err != nil {
			return api.Error(ctx, errs.
				New("Parse gateway request body").
				SetError(err))
		}

		gatewayRequest := NewRequest(requestBody, api.Headers(ctx), api.Cookies(ctx))
		gatewayResponse, err := matchService.Proxy(api.Context(ctx), matchRoute, gatewayRequest)
		if err != nil {
			return api.Error(ctx, err)
		}

		return api.SuccessRaw(ctx, gatewayResponse.StatusCode(), gatewayResponse.Body(), gatewayResponse.ContentType())
	}
}
