package http

import (
	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

type userHandler struct {
	userUsecase domain.UserUsecase
	logger      *logrus.Logger
}

func NewUserHandler(
	e *echo.Echo,
	us domain.UserUsecase,
	logger *logrus.Logger,
) {
	handler := &userHandler{
		userUsecase: us,
		logger:      logger,
	}
	e.POST("/sign-up", handler.SignUp)
	e.POST("/sign-in", handler.SignIn)
	e.POST("/sign-out", handler.SignOut)
	e.GET("/ping", handler.Ping)
	e.GET("/profile", handler.SeeProfileDetails)
	e.POST("/follow", handler.FollowUser)
	e.DELETE("/unfollow", handler.UnfollowUser)
	e.GET("/:user_id/following/list", handler.GetListOfFollowingHandler)
}
