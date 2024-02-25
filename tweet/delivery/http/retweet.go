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

func (handler *tweetHandler) Retweet(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.Retweet",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.Retweet",
		"request_id": requestId,
	}

	ctx = context.WithValue(context.Background(), "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.RetweetRequestPayload
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

	retweetUsecaseParam := domain.RetweetUsecaseParam{
		Token:   token,
		TweetId: reqPayload.TweetId,
	}

	logData["retweet_usecase_param"] = fmt.Sprintf("%+v", retweetUsecaseParam)

	retweetUsecaseResult, errRetweetUsecase := handler.tweetUsecase.Retweet(ctx, retweetUsecaseParam)
	if errRetweetUsecase != nil {
		logData["error_retweet"] = errRetweetUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errRetweetUsecase).
			Errorln("error call usecase Retweet")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errRetweetUsecase.Error())
	}

	logData["retweet_usecase_result"] = fmt.Sprintf("%+v", retweetUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success retweet")

	span.End()

	return echoCtx.JSON(http.StatusOK, retweetUsecaseResult)
}
