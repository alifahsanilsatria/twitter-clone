package http

import (
	"errors"
	"net/http"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/labstack/echo"

	"github.com/sirupsen/logrus"
)

func (handler *userHandler) SignOut(c echo.Context) error {
	logData := logrus.Fields{
		"method": "userHandler.SignOut",
	}

	token := c.Request().Header.Get("Token")

	if token == "" {
		return c.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	ctx := c.Request().Context()

	signOutUsecaseParam := domain.SignOutParam{
		Token: token,
	}
	signOutUsecaseResult, errorSignOutUsecase := handler.userUsecase.SignOut(ctx, signOutUsecaseParam)
	if errorSignOutUsecase != nil {
		logData["error_sign_up_usecase"] = errorSignOutUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errorSignOutUsecase).
			Errorln("error when parsing request payload")
		return c.JSON(http.StatusInternalServerError, errorSignOutUsecase.Error())
	}

	return c.JSON(http.StatusOK, signOutUsecaseResult)
}
