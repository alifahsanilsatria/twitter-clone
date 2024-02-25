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

func (handler *tweetHandler) PublishTweet(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.PublishTweet",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.PublishTweet",
		"request_id": requestId,
	}

	ctx = context.WithValue(context.Background(), "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.PublishTweetRequestPayload
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

	publishTweetUsecaseParam := domain.PublishTweetUsecaseParam{
		Token:    token,
		ParentId: reqPayload.ParentId,
		Content:  reqPayload.Content,
	}

	logData["publish_tweet_usecase_param"] = fmt.Sprintf("%+v", publishTweetUsecaseParam)

	publishTweetUsecaseResult, errPublishTweetUsecase := handler.tweetUsecase.PublishTweet(ctx, publishTweetUsecaseParam)
	if errPublishTweetUsecase != nil {
		logData["error_publish_tweet"] = errPublishTweetUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errPublishTweetUsecase).
			Errorln("error call usecase PublishTweet")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errPublishTweetUsecase.Error())
	}

	logData["publish_tweet_usecase_result"] = fmt.Sprintf("%+v", publishTweetUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success publish tweet")

	span.End()

	return echoCtx.JSON(http.StatusOK, publishTweetUsecaseResult)
}
