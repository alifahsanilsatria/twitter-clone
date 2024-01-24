package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/labstack/echo"

	"github.com/sirupsen/logrus"
)

func (handler *userHandler) GetListOfFollowersHandler(echoCtx echo.Context) error {
	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "userHandler.GetListOfFollowersHandler",
		"request_id": requestId,
	}

	ctx := context.WithValue(context.Background(), "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	userIdStr := echoCtx.Param("user_id")
	userId, errParseInt := strconv.ParseInt(userIdStr, 10, 32)
	if errParseInt != nil {
		return echoCtx.JSON(http.StatusBadRequest, errors.New("invalid tweet id"))
	}

	getListOfFollowersUsecaseParam := domain.GetListOfFollowersUsecaseParam{
		Token:  token,
		UserId: int32(userId),
	}

	logData["get_list_of_followers_usecase_param"] = fmt.Sprintf("%+v", getListOfFollowersUsecaseParam)

	getListOfFollowersUsecaseResult, errListOfFollowersUsecase := handler.userUsecase.GetListOfFollowersUsecase(ctx, getListOfFollowersUsecaseParam)
	if errListOfFollowersUsecase != nil {
		logData["error_get_list_of_followers_usecase"] = errListOfFollowersUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errListOfFollowersUsecase).
			Errorln("error on GetListOfFollowersUsecase")
		return echoCtx.JSON(http.StatusInternalServerError, errListOfFollowersUsecase.Error())
	}

	logData["get_list_of_followers_usecase_result"] = fmt.Sprintf("%+v", getListOfFollowersUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success get list of followers")

	return echoCtx.JSON(http.StatusOK, getListOfFollowersUsecaseResult)
}
