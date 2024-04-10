package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

// ExtraHeader will handle the ExtraHeader middleware
func (m *GoMiddleware) ExtraHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestId, _ := uuid.NewRandom()
		urlPath := ctx.Request().URL.Path

		ctx.Request().Header.Set("Request-Id", requestId.String())
		ctx.Request().Header.Set("Url-Path", urlPath)

		return next(ctx)
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
