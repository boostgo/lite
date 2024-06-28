package api

import (
	"github.com/boostgo/lite/errs"
	"github.com/labstack/echo/v4"
	"net/http"
)

func Parse(ctx echo.Context, export any) error {
	if err := ctx.Bind(export); err != nil {
		return errs.
			New("Parse request body error").
			SetError(err).
			SetHttpCode(http.StatusBadRequest)
	}

	// todo: validate
	return nil
}
