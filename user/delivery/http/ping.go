package http

import (
	"net/http"

	"github.com/labstack/echo"
)

func (handler *userHandler) Ping(c echo.Context) error {
	return c.JSON(http.StatusOK, "pong")
}
