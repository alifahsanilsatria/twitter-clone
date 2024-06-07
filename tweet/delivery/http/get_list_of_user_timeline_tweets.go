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

func (handler *tweetHandler) GetListOfUserTimelineTweets(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.GetListOfUserTimelineTweets",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.GetListOfUserTimelineTweets",
		"request_id": requestId,
	}

	ctx = context.WithValue(ctx, "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	getListOfUserTimelineTweetsParam := domain.GetListOfUserTimelineTweetsParam{
		Token: token,
	}

	logData["get_list_of_user_timeline_tweets_param"] = fmt.Sprintf("%+v", getListOfUserTimelineTweetsParam)

	getListOfUserTimelineTweetsResult, errGetListOfUserTimelineTweets := handler.tweetUsecase.GetListOfUserTimelineTweets(ctx, getListOfUserTimelineTweetsParam)
	if errGetListOfUserTimelineTweets != nil {
		logData["error_get_list_of_user_timeline_tweets"] = errGetListOfUserTimelineTweets.Error()
		handler.logger.
			WithFields(logData).
			WithError(errGetListOfUserTimelineTweets).
			Errorln("error call usecase GetListOfUserTimelineTweets")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errGetListOfUserTimelineTweets.Error())
	}

	logData["get_list_of_user_timeline_tweets_usecase_result"] = fmt.Sprintf("%+v", getListOfUserTimelineTweetsResult)
	handler.logger.
		WithFields(logData).
		Infoln("success get list of user timeline tweets")
	span.End()

	return echoCtx.JSON(http.StatusOK, getListOfUserTimelineTweetsResult)

}
