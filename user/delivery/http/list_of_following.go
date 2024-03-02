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
	"go.opentelemetry.io/otel/trace"
)

func (handler *userHandler) GetListOfFollowingHandler(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.GetListOfFollowingHandler",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "userHandler.GetListOfFollowingHandler",
		"request_id": requestId,
	}

	ctx = context.WithValue(ctx, "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	userIdStr := echoCtx.Param("user_id")
	userId, errParseInt := strconv.ParseInt(userIdStr, 10, 32)
	if errParseInt != nil {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("invalid tweet id"))
	}

	getListOfFollowingUsecaseParam := domain.GetListOfFollowingUsecaseParam{
		Token:  token,
		UserId: int32(userId),
	}

	logData["get_list_of_following_usecase_param"] = fmt.Sprintf("%+v", getListOfFollowingUsecaseParam)

	getListOfFollowingUsecaseResult, errListOfFollowingUsecase := handler.userUsecase.GetListOfFollowingUsecase(ctx, getListOfFollowingUsecaseParam)
	if errListOfFollowingUsecase != nil {
		logData["error_get_list_of_following_usecase"] = errListOfFollowingUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errListOfFollowingUsecase).
			Errorln("error on GetListOfFollowingUsecase")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errListOfFollowingUsecase.Error())
	}

	logData["get_list_of_following_usecase_result"] = fmt.Sprintf("%+v", getListOfFollowingUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success get list of following")
	span.End()

	return echoCtx.JSON(http.StatusOK, getListOfFollowingUsecaseResult)
}
