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

func (handler *tweetHandler) UndoLikes(echoCtx echo.Context) error {
	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "tweetHandler.UndoLikes",
		"request_id": requestId,
	}

	ctx := context.WithValue(context.Background(), "request_id", requestId)

	token := echoCtx.Request().Header.Get("Token")

	if token == "" {
		return echoCtx.JSON(http.StatusBadRequest, errors.New("empty token"))
	}

	logData["token"] = token

	var reqPayload domain.UndoLikesRequestPayload
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

	undoLikesUsecaseParam := domain.UndoLikesUsecaseParam{
		Token:   token,
		TweetId: reqPayload.TweetId,
	}

	logData["undo_likes_usecase_param"] = fmt.Sprintf("%+v", undoLikesUsecaseParam)

	undoLikesUsecaseResult, errUndoLikesUsecase := handler.tweetUsecase.UndoLikes(ctx, undoLikesUsecaseParam)
	if errUndoLikesUsecase != nil {
		logData["error_undo_likes"] = errUndoLikesUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errUndoLikesUsecase).
			Errorln("error call usecase UndoLikes")
		return echoCtx.JSON(http.StatusInternalServerError, errUndoLikesUsecase.Error())
	}

	logData["undo_likes_usecase_result"] = fmt.Sprintf("%+v", undoLikesUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success undo likes")

	return echoCtx.JSON(http.StatusOK, undoLikesUsecaseResult)
}
