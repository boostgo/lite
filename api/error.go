package api

import (
	"errors"
	"net/http"

	"github.com/boostgo/errorx"
)

// errStatusCode - define which status code must be provided to response by error
func errStatusCode(err error) int {
	switch {
	case errors.Is(err, errorx.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, errorx.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, errorx.ErrPaymentRequired):
		return http.StatusPaymentRequired
	case errors.Is(err, errorx.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, errorx.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, errorx.ErrMethodNotAllowed):
		return http.StatusMethodNotAllowed
	case errors.Is(err, errorx.ErrNotAcceptable):
		return http.StatusNotAcceptable
	case errors.Is(err, errorx.ErrProxyAuthRequired):
		return http.StatusProxyAuthRequired
	case errors.Is(err, errorx.ErrTimeout):
		return http.StatusRequestTimeout
	case errors.Is(err, errorx.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, errorx.ErrGone):
		return http.StatusGone
	case errors.Is(err, errorx.ErrLengthRequired):
		return http.StatusLengthRequired
	case errors.Is(err, errorx.ErrPreconditionFailed):
		return http.StatusPreconditionFailed
	case errors.Is(err, errorx.ErrEntityTooLarge):
		return http.StatusRequestEntityTooLarge
	case errors.Is(err, errorx.ErrURITooLong):
		return http.StatusRequestURITooLong
	case errors.Is(err, errorx.ErrUnsupportedMediaType):
		return http.StatusUnsupportedMediaType
	case errors.Is(err, errorx.ErrRangeNotSatisfiable):
		return http.StatusRequestedRangeNotSatisfiable
	case errors.Is(err, errorx.ErrExpectationFailed):
		return http.StatusExpectationFailed
	case errors.Is(err, errorx.ErrTeapot):
		return http.StatusTeapot
	case errors.Is(err, errorx.ErrMisdirectedRequest):
		return http.StatusMisdirectedRequest
	case errors.Is(err, errorx.ErrUnprocessableEntity):
		return http.StatusUnprocessableEntity
	case errors.Is(err, errorx.ErrLocked):
		return http.StatusLocked
	case errors.Is(err, errorx.ErrFailedDependency):
		return http.StatusFailedDependency
	case errors.Is(err, errorx.ErrTooEarly):
		return http.StatusTooEarly
	case errors.Is(err, errorx.ErrUpgradeRequired):
		return http.StatusUpgradeRequired
	case errors.Is(err, errorx.ErrPreconditionRequired):
		return http.StatusPreconditionRequired
	case errors.Is(err, errorx.ErrTooManyRequests):
		return http.StatusTooManyRequests
	case errors.Is(err, errorx.ErrRequestHeaderFieldsTooLarge):
		return http.StatusRequestHeaderFieldsTooLarge
	case errors.Is(err, errorx.ErrUnavailableForLegalReasons):
		return http.StatusUnavailableForLegalReasons
	case errors.Is(err, errorx.ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, errorx.ErrNotImplemented):
		return http.StatusNotImplemented
	case errors.Is(err, errorx.ErrBadGateway):
		return http.StatusBadGateway
	case errors.Is(err, errorx.ErrServiceUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, errorx.ErrGatewayTimeout):
		return http.StatusGatewayTimeout
	case errors.Is(err, errorx.ErrHTTPVersionNotSupported):
		return http.StatusHTTPVersionNotSupported
	case errors.Is(err, errorx.ErrVariantAlsoNegotiates):
		return http.StatusVariantAlsoNegotiates
	case errors.Is(err, errorx.ErrInsufficientStorage):
		return http.StatusInsufficientStorage
	case errors.Is(err, errorx.ErrLoopDetected):
		return http.StatusLoopDetected
	case errors.Is(err, errorx.ErrNotExtended):
		return http.StatusNotExtended
	case errors.Is(err, errorx.ErrNetworkAuthenticationFailed):
		return http.StatusNetworkAuthenticationRequired
	default:
		return http.StatusInternalServerError
	}
}
