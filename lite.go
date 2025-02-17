package lite

import (
	"errors"
	"github.com/boostgo/lite/app/api"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/system/try"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	handler *echo.Echo
	_debug  = atomic.Bool{}
	_podID  string
)

func init() {
	_podID = uuid.New().String()

	handler = echo.New()

	handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Auth-Token"},
		AllowCredentials: true,
	}))
	handler.Use(RecoverMiddleware())

	if trace.AmIMaster() {
		handler.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {
				trace.SetEchoCtx(ctx, trace.String())
				return next(ctx)
			}
		})
	}

	handler.RouteNotFound("*", func(ctx echo.Context) error {
		return api.Error(ctx, errs.
			New("Route not found").
			SetError(errs.ErrNotFound).
			AddContext("url", ctx.Request().RequestURI))
	})
}

func TimeoutMiddleware(duration time.Duration) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request reached timeout",
		OnTimeoutRouteErrorHandler: func(err error, ctx echo.Context) {
			_ = api.Error(
				ctx,
				errs.
					New("Request reached timeout").
					SetError(err, errs.ErrTimeout),
			)
		},
		Timeout: duration,
	})
}

func RecoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if err := try.Try(func() error {
				return next(ctx)
			}); err != nil {
				return api.Error(ctx, err)
			}

			return nil
		}
	}
}

func SetDebug(debug bool) {
	_debug.Store(debug)
	handler.Debug = debug
}

func With(fn func(h *echo.Echo)) {
	fn(handler)
}

func Debug() bool {
	return _debug.Load()
}

func Use(middlewares ...echo.MiddlewareFunc) {
	handler.Use(middlewares...)
}

func run(address string) error {
	if trace.AmIMaster() {
		handler.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator:        uuid.NewString,
			RequestIDHandler: trace.SetEchoCtx,
			TargetHeader:     trace.Key(),
		}))
	}

	life.Tear(func() error {
		return handler.Shutdown(life.Context())
	})

	if err := handler.Start(address); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return errs.New("Start server").SetError(err)
	}

	return nil
}

func Handler() *echo.Echo {
	return handler
}

func Run(address string, waitTime ...time.Duration) {
	go func() {
		if err := run(address); err != nil {
			log.
				Namespace("handler").
				Error().
				Err(err).
				Send()
			life.Cancel()
		}
	}()

	life.GracefulLog(func() {
		log.
			Namespace("lite").
			Info().
			Msg("Graceful shutdown...")
	})
	life.Wait(waitTime...)
}

func PodID() string {
	return _podID
}

type Router interface {
	Any(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) []*echo.Route
	POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Use(m ...echo.MiddlewareFunc)
	RouteNotFound(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route
	Group(prefix string, m ...echo.MiddlewareFunc) (g *echo.Group)
}
