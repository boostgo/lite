package lite

import (
	"context"
	"errors"
	"net/http"
	"sync/atomic"
	"time"

	"github.com/boostgo/appx"
	"github.com/boostgo/errorx"
	"github.com/boostgo/lite/api"
	"github.com/boostgo/log"
	"github.com/boostgo/trace"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	handler *echo.Echo
	_debug  = atomic.Bool{}
	_podID  string
	_tracer *trace.Tracer
)

func InitTracer(tracer *trace.Tracer) {
	_tracer = tracer
}

func init() {
	_podID = uuid.New().String()

	handler = echo.New()

	handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Auth-Token"},
		AllowCredentials: true,
	}))
	handler.Use(RecoverMiddleware())

	handler.RouteNotFound("*", func(ctx echo.Context) error {
		return api.Error(ctx, errorx.
			New("Route not found").
			SetError(errorx.ErrNotFound).
			AddContext("url", ctx.Request().RequestURI))
	})
}

func TimeoutMiddleware(duration time.Duration) echo.MiddlewareFunc {
	return middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{
		Skipper: middleware.DefaultSkipper,
		ErrorHandler: func(err error, ctx echo.Context) error {
			return api.Error(
				ctx,
				errorx.
					New("Request reached timeout").
					SetError(err, errorx.ErrTimeout),
			)
		},
		Timeout: duration,
	})
}

func RecoverMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			if err := errorx.Try(func() error {
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
	if _tracer.AmIMaster() {
		handler.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
			Generator: uuid.NewString,
			RequestIDHandler: func(ctx echo.Context, traceID string) {
				ctx.SetRequest(
					ctx.Request().WithContext(
						context.WithValue(
							ctx.Request().Context(),
							"bgo_trace_id",
							traceID,
						)))
			},
			TargetHeader: "X-Trace-ID",
		}))
	}

	appx.Tear(func() error {
		return handler.Shutdown(appx.Context())
	})

	if err := handler.Start(address); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}

		return errorx.New("Start server").SetError(err)
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
				Error().
				Err(err).
				Msg("Run server")
			appx.Cancel()
		}
	}()

	appx.GracefulLog(func() {
		log.
			Info().
			Msg("Graceful shutdown...")
	})
	appx.Wait(waitTime...)
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
