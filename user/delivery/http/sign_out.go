package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func (handler *userHandler) SignOut(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.SignOut",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "userHandler.SignOut",
		"request_id": requestId,
	}

	ctx = context.WithValue(context.Background(), "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	signOutUsecaseParam := domain.SignOutUsecaseParam{
		Token: token,
	}

	logData["sign_out_usecase_param"] = fmt.Sprintf("%+v", signOutUsecaseParam)

	signOutUsecaseResult, errorSignOutUsecase := handler.userUsecase.SignOut(ctx, signOutUsecaseParam)
	if errorSignOutUsecase != nil {
		logData["error_sign_out_usecase"] = errorSignOutUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errorSignOutUsecase).
			Errorln("error when parsing request payload")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errorSignOutUsecase.Error())
	}

	logData["sign_out_usecase_result"] = fmt.Sprintf("%+v", signOutUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success sign out")
	span.End()

	return echoCtx.JSON(http.StatusOK, signOutUsecaseResult)
}
