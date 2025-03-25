package api

import (
	"errors"
	"github.com/boostgo/lite/errs"
	"net/http"
)

// errStatusCode - define which status code must be provided to response by error
func errStatusCode(err error) int {
	switch {
	case errors.Is(err, errs.ErrBadRequest):
		return http.StatusBadRequest
	case errors.Is(err, errs.ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, errs.ErrPaymentRequired):
		return http.StatusPaymentRequired
	case errors.Is(err, errs.ErrForbidden):
		return http.StatusForbidden
	case errors.Is(err, errs.ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, errs.ErrMethodNotAllowed):
		return http.StatusMethodNotAllowed
	case errors.Is(err, errs.ErrNotAcceptable):
		return http.StatusNotAcceptable
	case errors.Is(err, errs.ErrProxyAuthRequired):
		return http.StatusProxyAuthRequired
	case errors.Is(err, errs.ErrTimeout):
		return http.StatusRequestTimeout
	case errors.Is(err, errs.ErrConflict):
		return http.StatusConflict
	case errors.Is(err, errs.ErrGone):
		return http.StatusGone
	case errors.Is(err, errs.ErrLengthRequired):
		return http.StatusLengthRequired
	case errors.Is(err, errs.ErrPreconditionFailed):
		return http.StatusPreconditionFailed
	case errors.Is(err, errs.ErrEntityTooLarge):
		return http.StatusRequestEntityTooLarge
	case errors.Is(err, errs.ErrURITooLong):
		return http.StatusRequestURITooLong
	case errors.Is(err, errs.ErrUnsupportedMediaType):
		return http.StatusUnsupportedMediaType
	case errors.Is(err, errs.ErrRangeNotSatisfiable):
		return http.StatusRequestedRangeNotSatisfiable
	case errors.Is(err, errs.ErrExpectationFailed):
		return http.StatusExpectationFailed
	case errors.Is(err, errs.ErrTeapot):
		return http.StatusTeapot
	case errors.Is(err, errs.ErrMisdirectedRequest):
		return http.StatusMisdirectedRequest
	case errors.Is(err, errs.ErrUnprocessableEntity):
		return http.StatusUnprocessableEntity
	case errors.Is(err, errs.ErrLocked):
		return http.StatusLocked
	case errors.Is(err, errs.ErrFailedDependency):
		return http.StatusFailedDependency
	case errors.Is(err, errs.ErrTooEarly):
		return http.StatusTooEarly
	case errors.Is(err, errs.ErrUpgradeRequired):
		return http.StatusUpgradeRequired
	case errors.Is(err, errs.ErrPreconditionRequired):
		return http.StatusPreconditionRequired
	case errors.Is(err, errs.ErrTooManyRequests):
		return http.StatusTooManyRequests
	case errors.Is(err, errs.ErrRequestHeaderFieldsTooLarge):
		return http.StatusRequestHeaderFieldsTooLarge
	case errors.Is(err, errs.ErrUnavailableForLegalReasons):
		return http.StatusUnavailableForLegalReasons
	case errors.Is(err, errs.ErrInternal):
		return http.StatusInternalServerError
	case errors.Is(err, errs.ErrNotImplemented):
		return http.StatusNotImplemented
	case errors.Is(err, errs.ErrBadGateway):
		return http.StatusBadGateway
	case errors.Is(err, errs.ErrServiceUnavailable):
		return http.StatusServiceUnavailable
	case errors.Is(err, errs.ErrGatewayTimeout):
		return http.StatusGatewayTimeout
	case errors.Is(err, errs.ErrHTTPVersionNotSupported):
		return http.StatusHTTPVersionNotSupported
	case errors.Is(err, errs.ErrVariantAlsoNegotiates):
		return http.StatusVariantAlsoNegotiates
	case errors.Is(err, errs.ErrInsufficientStorage):
		return http.StatusInsufficientStorage
	case errors.Is(err, errs.ErrLoopDetected):
		return http.StatusLoopDetected
	case errors.Is(err, errs.ErrNotExtended):
		return http.StatusNotExtended
	case errors.Is(err, errs.ErrNetworkAuthenticationFailed):
		return http.StatusNetworkAuthenticationRequired
	default:
		return http.StatusInternalServerError
	}
}
