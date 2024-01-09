package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (handler *tweetHandler) DeleteTweet(c echo.Context) error {
	logData := logrus.Fields{
		"method": "tweetHandler.PublishTweet",
	}

	token := c.Request().Header.Get("Token")

	if token == "" {
		return c.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.DeleteTweetRequestPayload
	errParsingReqPayload := json.NewDecoder(c.Request().Body).Decode(&reqPayload)
	if errParsingReqPayload != nil {
		logData["error_parsing_request_payload"] = errParsingReqPayload.Error()
		handler.logger.
			WithFields(logData).
			WithError(errParsingReqPayload).
			Errorln("error when parsing request payload")
		return c.JSON(http.StatusUnprocessableEntity, errParsingReqPayload.Error())
	}

	ctx := c.Request().Context()

	deleteTweetParam := domain.DeleteTweetParam{
		Token:   token,
		TweetId: reqPayload.TweetId,
		UserId:  reqPayload.UserId,
	}
	deleteTweetUsecaseResult, errDeleteTweetUsecase := handler.tweetUsecase.DeleteTweet(ctx, deleteTweetParam)
	if errDeleteTweetUsecase != nil {
		logData["error_delete_tweet"] = errDeleteTweetUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errDeleteTweetUsecase).
			Errorln("error call usecase DeleteTweet")
		return c.JSON(http.StatusInternalServerError, errDeleteTweetUsecase.Error())
	}

	return c.JSON(http.StatusOK, deleteTweetUsecaseResult)
}
