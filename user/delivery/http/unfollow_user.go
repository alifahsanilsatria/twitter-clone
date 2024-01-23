package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (handler *userHandler) UnfollowUser(echoCtx echo.Context) error {
	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "userHandler.UnfollowUser",
		"request_id": echoCtx.Request().Header.Get("Request-Id"),
	}

	ctx := context.WithValue(context.Background(), "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.FollowUserRequestPayload
	errParsingReqPayload := json.NewDecoder(echoCtx.Request().Body).Decode(&reqPayload)
	if errParsingReqPayload != nil {
		logData["error_parsing_request_payload"] = errParsingReqPayload.Error()
		handler.logger.
			WithFields(logData).
			WithError(errParsingReqPayload).
			Errorln("error when parsing request payload")
		return echoCtx.JSON(http.StatusUnprocessableEntity, errParsingReqPayload.Error())
	}

	logData["request_payload"] = fmt.Sprintf("%+v", reqPayload)

	unfollowUserUsecaseParam := domain.UnfollowUserParam{
		Token:           token,
		FollowingUserId: reqPayload.FollowingUserId,
	}

	logData["unfollow_user_usecase_param"] = fmt.Sprintf("%+v", unfollowUserUsecaseParam)

	unfollowUserUsecaseResult, errorUnfollowUserUsecase := handler.userUsecase.UnfollowUser(ctx, unfollowUserUsecaseParam)
	if errorUnfollowUserUsecase != nil {
		logData["error_unfollow_user_usecase"] = errorUnfollowUserUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errorUnfollowUserUsecase).
			Errorln("error when unfollow user")
		return echoCtx.JSON(http.StatusInternalServerError, errorUnfollowUserUsecase.Error())
	}

	logData["unfollow_user_usecase_result"] = fmt.Sprintf("%+v", unfollowUserUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success unfollow user")

	return echoCtx.JSON(http.StatusOK, unfollowUserUsecaseResult)
}
