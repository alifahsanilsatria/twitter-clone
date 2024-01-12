package middleware

import (
	"github.com/google/uuid"
	"github.com/labstack/echo"
)

// GoMiddleware represent the data-struct for middleware
type GoMiddleware struct {
	// another stuff , may be needed by middleware
}

// RequestId will handle the RequestId middleware
func (m *GoMiddleware) RequestId(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		requestId, _ := uuid.NewRandom()
		ctx.Request().Header.Set("Request-Id", requestId.String())
		return next(ctx)
	}
}

// InitMiddleware initialize the middleware
func InitMiddleware() *GoMiddleware {
	return &GoMiddleware{}
}
