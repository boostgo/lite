package api

import (
	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func Health(router *echo.Echo) {
	router.GET("/health", func(ctx echo.Context) error {
		return Ok(ctx, "OK")
	})
}

func Swagger(router *echo.Echo) {
	router.GET("/swagger/*", echoSwagger.WrapHandler)
}
