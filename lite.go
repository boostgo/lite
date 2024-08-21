package lite

import (
	"errors"
	"github.com/boostgo/lite/app/api"
	"github.com/boostgo/lite/errs"
	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/system/life"
	"github.com/boostgo/lite/system/trace"
	"github.com/boostgo/lite/system/try"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	handler *echo.Echo
	_debug  = atomic.Bool{}
)

func init() {
	handler = echo.New()

	handler.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
	}))
	handler.Use(RecoverMiddleware())
	handler.Use(middleware.RequestID())
	handler.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request reached timeout",
		OnTimeoutRouteErrorHandler: func(err error, ctx echo.Context) {
			_ = api.Error(ctx, errs.New("Request reached timeout").SetError(err, errs.ErrTimeout))
		},
		Timeout: time.Second * 30,
	}))

	handler.RouteNotFound("*", func(ctx echo.Context) error {
		return api.Error(ctx, errs.New("Route not found").SetError(errs.ErrNotFound))
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
		handler.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(ctx echo.Context) error {
				trace.SetEchoCtx(ctx, trace.String())
				return next(ctx)
			}
		})
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

func Run(address string) {
	go func() {
		if err := run(address); err != nil {
			log.Error().Err(err).Namespace("handler")
			life.Cancel()
		}
	}()

	life.GracefulLog(func() {
		log.Info().Msg("Graceful shutdown...").Namespace("lite")
	})
	life.Wait()
}
