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

func (handler *userHandler) SeeProfileDetails(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.SeeProfileDetails",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "userHandler.SeeProfileDetails",
		"request_id": requestId,
	}

	ctx = context.WithValue(ctx, "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	seeProfileDetailsUsecaseParam := domain.SeeProfileDetailsParam{
		Token: token,
	}

	logData["see_profile_details_usecase_param"] = fmt.Sprintf("%+v", seeProfileDetailsUsecaseParam)

	seeProfileDetailsUsecaseResult, errorSeeProfileDetailsUsecase := handler.userUsecase.SeeProfileDetails(ctx, seeProfileDetailsUsecaseParam)
	if errorSeeProfileDetailsUsecase != nil {
		logData["error_see_profile_details_usecase"] = errorSeeProfileDetailsUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errorSeeProfileDetailsUsecase).
			Errorln("error on SeeProfileDetails")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errorSeeProfileDetailsUsecase.Error())
	}

	logData["see_profile_details_usecase_result"] = fmt.Sprintf("%+v", seeProfileDetailsUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success see profile details")
	span.End()

	return echoCtx.JSON(http.StatusOK, seeProfileDetailsUsecaseResult)
}
