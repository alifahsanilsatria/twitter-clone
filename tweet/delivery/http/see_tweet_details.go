package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/trace"
)

func (handler *tweetHandler) SeeTweetDetails(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.SeeTweetDetails",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.SeeTweetDetails",
		"request_id": requestId,
	}

	ctx = context.WithValue(ctx, "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	tweetIdStr := echoCtx.Param("tweet_id")
	tweetId, errParseInt := strconv.ParseInt(tweetIdStr, 10, 32)
	if errParseInt != nil {
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errors.New("invalid tweet id"))
	}

	logData["token"] = token
	logData["tweet_id"] = fmt.Sprintf("%+v", tweetIdStr)

	seeTweetDetailsUsecaseParam := domain.SeeTweetDetailsUsecaseParam{
		Token:   token,
		TweetId: int32(tweetId),
	}

	logData["see_tweet_details_usecase_param"] = fmt.Sprintf("%+v", seeTweetDetailsUsecaseParam)

	seeTweetDetailsUsecaseResult, errSeeTweetDetails := handler.tweetUsecase.SeeTweetDetails(ctx, seeTweetDetailsUsecaseParam)
	if errSeeTweetDetails != nil {
		logData["error_see_tweet_details"] = errSeeTweetDetails.Error()
		handler.logger.
			WithFields(logData).
			WithError(errSeeTweetDetails).
			Errorln("error call usecase SeeTweetDetails")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errSeeTweetDetails.Error())
	}

	logData["see_tweet_details_usecase_result"] = fmt.Sprintf("%+v", seeTweetDetailsUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success see tweet details")
	span.End()

	return echoCtx.JSON(http.StatusOK, seeTweetDetailsUsecaseResult)

}
