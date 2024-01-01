package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (handler *tweetHandler) PublishTweet(c echo.Context) error {
	logData := logrus.Fields{
		"method": "tweetHandler.PublishTweet",
	}

	token := c.Request().Header.Get("Token")

	if token == "" {
		return c.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.PublishTweetRequestPayload
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

	publishTweetParam := domain.PublishTweetParam{
		Token:    token,
		ParentId: reqPayload.ParentId,
		Content:  reqPayload.Content,
	}
	publishTweetUsecaseResult, errPublishTweetUsecase := handler.tweetUsecase.PublishTweet(ctx, publishTweetParam)
	if errPublishTweetUsecase != nil {
		logData["error_publish_tweet"] = errPublishTweetUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errPublishTweetUsecase).
			Errorln("error call usecase PublishTweet")
		return c.JSON(http.StatusInternalServerError, errPublishTweetUsecase.Error())
	}

	return c.JSON(http.StatusOK, publishTweetUsecaseResult)

}
