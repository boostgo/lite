package lite

import "github.com/labstack/echo/v4"

func GET(path string, pathHandler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	handler.GET(path, pathHandler, middlewares...)
}

func POST(path string, pathHandler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	handler.POST(path, pathHandler, middlewares...)
}

func PUT(path string, pathHandler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	handler.PUT(path, pathHandler, middlewares...)
}

func DELETE(path string, pathHandler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	handler.DELETE(path, pathHandler, middlewares...)
}

func HEAD(path string, pathHandler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	handler.HEAD(path, pathHandler, middlewares...)
}

func OPTIONS(path string, pathHandler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	handler.OPTIONS(path, pathHandler, middlewares...)
}

func NotFound(path string, pathHandler echo.HandlerFunc, middlewares ...echo.MiddlewareFunc) {
	handler.RouteNotFound(path, pathHandler, middlewares...)
}

func Group(prefix string) *echo.Group {
	return handler.Group(prefix)
}
