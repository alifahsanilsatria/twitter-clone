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
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (handler *userHandler) FollowUser(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.FollowUser",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "userHandler.FollowUser",
		"request_id": echoCtx.Request().Header.Get("Request-Id"),
	}

	ctx = context.WithValue(ctx, "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
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
		span.End()
		return echoCtx.JSON(http.StatusUnprocessableEntity, errParsingReqPayload.Error())
	}

	span.SetAttributes(attribute.String("request_id", requestId))
	span.SetAttributes(attribute.String("request_payload", fmt.Sprintf("%+v", reqPayload)))

	logData["request_payload"] = fmt.Sprintf("%+v", reqPayload)

	followUserUsecaseParam := domain.FollowUserParam{
		Token:           token,
		FollowingUserId: reqPayload.FollowingUserId,
	}

	logData["follow_user_usecase_param"] = fmt.Sprintf("%+v", followUserUsecaseParam)

	followUserUsecaseResult, errorFollowUserUsecase := handler.userUsecase.FollowUser(ctx, followUserUsecaseParam)
	if errorFollowUserUsecase != nil {
		logData["error_follow_user_usecase"] = errorFollowUserUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errorFollowUserUsecase).
			Errorln("error when follow user")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errorFollowUserUsecase.Error())
	}

	logData["follow_user_usecase_result"] = fmt.Sprintf("%+v", followUserUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success follow user")
	span.End()

	return echoCtx.JSON(http.StatusOK, followUserUsecaseResult)
}
