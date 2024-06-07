package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func (handler *tweetHandler) GetListOfAUserTimelineTweets(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.GetListOfAUserTimelineTweets",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.GetListOfAUserTimelineTweets",
		"request_id": requestId,
	}

	ctx = context.WithValue(ctx, "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	username := echoCtx.Param("user_name")

	logData["token"] = token

	getListOfAUserTimelineTweetsParam := domain.GetListOfAUserTimelineTweetsParam{
		Token:    token,
		Username: username,
	}

	logData["get_list_of_user_timeline_tweets_param"] = fmt.Sprintf("%+v", getListOfAUserTimelineTweetsParam)

	getListOfAUserTimelineTweetsResult, errGetListOfAUserTimelineTweets := handler.tweetUsecase.GetListOfAUserTimelineTweets(ctx, getListOfAUserTimelineTweetsParam)
	if errGetListOfAUserTimelineTweets != nil {
		logData["error_get_list_of_a_user_timeline_tweets"] = errGetListOfAUserTimelineTweets.Error()
		handler.logger.
			WithFields(logData).
			WithError(errGetListOfAUserTimelineTweets).
			Errorln("error call usecase GetListOfAUserTimelineTweets")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errGetListOfAUserTimelineTweets.Error())
	}

	logData["get_list_of_a_user_timeline_tweets_usecase_result"] = fmt.Sprintf("%+v", getListOfAUserTimelineTweetsResult)
	handler.logger.
		WithFields(logData).
		Infoln("success get list of user timeline tweets")
	span.End()

	return echoCtx.JSON(http.StatusOK, getListOfAUserTimelineTweetsResult)

}
