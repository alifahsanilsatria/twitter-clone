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

func (handler *tweetHandler) UndoRetweet(echoCtx echo.Context) error {
	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.UndoRetweet",
		"request_id": requestId,
	}

	ctx := context.WithValue(context.Background(), "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.UndoRetweetRequestPayload
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

	undoRetweetUsecaseParam := domain.UndoRetweetUsecaseParam{
		Token:   token,
		TweetId: reqPayload.TweetId,
	}

	logData["undo_retweet_usecase_param"] = fmt.Sprintf("%+v", undoRetweetUsecaseParam)

	undoRetweetUsecaseResult, errUndoRetweetUsecase := handler.tweetUsecase.UndoRetweet(ctx, undoRetweetUsecaseParam)
	if errUndoRetweetUsecase != nil {
		logData["error_undo_retweet"] = errUndoRetweetUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errUndoRetweetUsecase).
			Errorln("error call usecase UndoRetweet")
		return echoCtx.JSON(http.StatusInternalServerError, errUndoRetweetUsecase.Error())
	}

	logData["undo_retweet_usecase_result"] = fmt.Sprintf("%+v", undoRetweetUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success undo retweet")

	return echoCtx.JSON(http.StatusOK, undoRetweetUsecaseResult)
}
