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
)

func (handler *tweetHandler) LikeTweet(echoCtx echo.Context) error {
	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.LikeTweet",
		"request_id": requestId,
	}

	ctx := context.WithValue(context.Background(), "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.LikeTweetRequestPayload
	errParsingReqPayload := json.NewDecoder(echoCtx.Request().Body).Decode(&reqPayload)
	if errParsingReqPayload != nil {
		logData["error_parsing_request_payload"] = errParsingReqPayload.Error()
		handler.logger.
			WithFields(logData).
			WithError(errParsingReqPayload).
			Errorln("error when parsing request payload")
		return echoCtx.JSON(http.StatusUnprocessableEntity, errParsingReqPayload.Error())
	}

	logData["request_payload"] = fmt.Sprintf("%+v", reqPayload)

	likeTweetUsecaseParam := domain.LikeTweetUsecaseParam{
		Token:   token,
		TweetId: reqPayload.TweetId,
	}

	logData["like_tweet_usecase_param"] = fmt.Sprintf("%+v", likeTweetUsecaseParam)

	likeTweetUsecaseResult, errLikeTweetUsecase := handler.tweetUsecase.LikeTweet(ctx, likeTweetUsecaseParam)
	if errLikeTweetUsecase != nil {
		logData["error_like_tweet"] = errLikeTweetUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errLikeTweetUsecase).
			Errorln("error call usecase LikeTweet")
		return echoCtx.JSON(http.StatusInternalServerError, errLikeTweetUsecase.Error())
	}

	logData["like_tweet_usecase_result"] = fmt.Sprintf("%+v", likeTweetUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success like tweet")

	return echoCtx.JSON(http.StatusOK, likeTweetUsecaseResult)
}
