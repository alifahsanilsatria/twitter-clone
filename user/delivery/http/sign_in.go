package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/common"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func (handler *userHandler) SignIn(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.SignIn",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	var pushToPrometheusMetricsParam = common.PushMetricsToPrometheusParam{
		Path: echoCtx.Request().Header.Get("Url-Path"),
	}

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "userHandler.SignIn",
		"request_id": requestId,
	}

	ctx = context.WithValue(ctx, "request_id", requestId)

	var reqPayload domain.SignInRequestPayload
	errParsingReqPayload := json.NewDecoder(echoCtx.Request().Body).Decode(&reqPayload)
	if errParsingReqPayload != nil {
		logData["error_parsing_request_payload"] = errParsingReqPayload.Error()
		handler.logger.
			WithFields(logData).
			WithError(errParsingReqPayload).
			Errorln("error when parsing request payload")
		span.End()
		common.PushMetricsToPrometheus(pushToPrometheusMetricsParam)
		return echoCtx.JSON(http.StatusUnprocessableEntity, errParsingReqPayload.Error())
	}

	logData["request_payload"] = fmt.Sprintf("%+v", reqPayload)

	errvalidateSignInParam := validateSignInParam(reqPayload)
	if errParsingReqPayload != nil {
		logData["error_validate_sign_in_param"] = errvalidateSignInParam.Error()
		handler.logger.
			WithFields(logData).
			WithError(errvalidateSignInParam).
			Errorln("error when validate sign in param")
		span.End()
		common.PushMetricsToPrometheus(pushToPrometheusMetricsParam)
		return echoCtx.JSON(http.StatusBadRequest, errvalidateSignInParam.Error())
	}

	signInUsecaseParam := domain.SignInUsecaseParam{
		Username: reqPayload.Username,
		Password: reqPayload.Password,
	}

	logData["sign_in_usecase_param"] = fmt.Sprintf("%+v", signInUsecaseParam)

	signInUsecaseResult, errorSignInUsecase := handler.userUsecase.SignIn(ctx, signInUsecaseParam)
	if errorSignInUsecase != nil {
		logData["error_sign_in_usecase"] = errorSignInUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errorSignInUsecase).
			Errorln("error when sign in")
		span.End()
		common.PushMetricsToPrometheus(pushToPrometheusMetricsParam)
		return echoCtx.JSON(http.StatusInternalServerError, errorSignInUsecase.Error())
	}

	logData["sign_in_usecase_result"] = fmt.Sprintf("%+v", signInUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success sign in")
	span.End()
	pushToPrometheusMetricsParam.IsRequestSuccessful = true
	common.PushMetricsToPrometheus(pushToPrometheusMetricsParam)

	return echoCtx.JSON(http.StatusOK, signInUsecaseResult)
}

func validateSignInParam(param domain.SignInRequestPayload) error {
	if param.Username == "" {
		return errors.New("username is empty")
	}

	if param.Password == "" {
		return errors.New("password is empty")
	}

	return nil
}
