package api

import (
	"bytes"
	"context"
	"github.com/boostgo/convert"
	"io"
	"net/http"
	"time"

	"github.com/boostgo/lite/log"
	"github.com/boostgo/lite/types/content"
	"github.com/labstack/echo/v4"
)

const (
	rawResponseKey = "lite-response-raw"
)

// Raw if middleware set, all responses by this middleware will be returned in "raw" way (no successOutput object)
func Raw() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			localCtx := Context(ctx)
			localCtx = context.WithValue(localCtx, rawResponseKey, true)
			ctx.SetRequest(ctx.Request().WithContext(localCtx))
			return next(ctx)
		}
	}
}

func Cache(ttl time.Duration, distributor HttpCacheDistributor) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// try load response from cache
			responseBody, cacheOk, err := distributor.Get(Context(ctx), ctx.Request())
			if err != nil {
				cacheOk = false

				log.
					Error(Context(ctx)).
					Err(err).
					Msg("Get response body by http cache distributor")
			}

			// return cached response
			if cacheOk {
				return SuccessRaw(ctx, http.StatusOK, responseBody, content.JSON)
			}

			// call handler method to generate response
			response := ctx.Response()
			var responseBuffer bytes.Buffer
			mw := io.MultiWriter(&responseBuffer, response.Writer)
			response.Writer = newCacheResponseWriter(response.Writer, mw)

			if err = next(ctx); err != nil {
				return err
			}

			responseBody = responseBuffer.Bytes()

			// set response to cache
			if err = distributor.Set(Context(ctx), ctx.Request(), responseBody, ttl); err != nil {
				log.
					Error(Context(ctx)).
					Err(err).
					Msg("Set response body by http cache distributor")
			}

			return nil
		}
	}
}

func isRaw(ctx echo.Context) bool {
	return convert.Bool(Context(ctx).Value(rawResponseKey))
}
