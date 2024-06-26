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

func (handler *tweetHandler) DeleteTweet(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.DeleteTweet",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.DeleteTweet",
		"request_id": requestId,
	}

	ctx = context.WithValue(ctx, "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.DeleteTweetRequestPayload
	errParsingReqPayload := json.NewDecoder(echoCtx.Request().Body).Decode(&reqPayload)
	if errParsingReqPayload != nil {
		logData["error_parsing_request_payload"] = errParsingReqPayload.Error()
		handler.logger.
			WithFields(logData).
			WithError(errParsingReqPayload).
			Errorln("error when parsing request payload")
		return echoCtx.JSON(http.StatusUnprocessableEntity, errParsingReqPayload.Error())
	}

	span.SetAttributes(attribute.String("request_id", requestId))
	span.SetAttributes(attribute.String("request_payload", fmt.Sprintf("%+v", reqPayload)))

	logData["request_payload"] = fmt.Sprintf("%+v", reqPayload)

	deleteTweetUsecaseParam := domain.DeleteTweetUsecaseParam{
		Token:   token,
		TweetId: reqPayload.TweetId,
	}

	logData["delete_tweet_usecase_param"] = fmt.Sprintf("%+v", deleteTweetUsecaseParam)

	deleteTweetUsecaseResult, errDeleteTweetUsecase := handler.tweetUsecase.DeleteTweet(ctx, deleteTweetUsecaseParam)
	if errDeleteTweetUsecase != nil {
		logData["error_delete_tweet"] = errDeleteTweetUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errDeleteTweetUsecase).
			Errorln("error call usecase DeleteTweet")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errDeleteTweetUsecase.Error())
	}

	logData["delete_tweet_usecase_result"] = fmt.Sprintf("%+v", deleteTweetUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success delete tweet")

	span.End()

	return echoCtx.JSON(http.StatusOK, deleteTweetUsecaseResult)
}
