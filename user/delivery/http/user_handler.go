package http

import (
	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	trace "go.opentelemetry.io/otel/trace"
)

type userHandler struct {
	userUsecase domain.UserUsecase
	logger      *logrus.Logger
	tracer      trace.Tracer
}

func NewUserHandler(
	e *echo.Echo,
	us domain.UserUsecase,
	logger *logrus.Logger,
	tracer trace.Tracer,
) {
	handler := &userHandler{
		userUsecase: us,
		logger:      logger,
		tracer:      tracer,
	}
	e.POST("/sign-up", handler.SignUp)
	e.POST("/sign-in", handler.SignIn)
	e.POST("/sign-out", handler.SignOut)
	e.GET("/ping", handler.Ping)
}
