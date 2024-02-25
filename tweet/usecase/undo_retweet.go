package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (uc *tweetUsecase) UndoRetweet(ctx context.Context, param domain.UndoRetweetUsecaseParam) (domain.UndoRetweetUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.UndoRetweet", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.UndoRetweet",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	getUserSessionByTokenParam := domain.GetUserSessionByTokenParam{
		Token: param.Token,
	}

	logData["get_user_session_by_token_param"] = fmt.Sprintf("%+v", getUserSessionByTokenParam)

	userSession, errGetUserSession := uc.userSessionRepository.GetUserSessionByToken(ctx, getUserSessionByTokenParam)
	if errGetUserSession != nil {
		logData["error_get_user_session"] = errGetUserSession.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserSession).
			Errorln("error on GetUserSessionByToken")
		span.End()
		return domain.UndoRetweetUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.UndoRetweetUsecaseResult{}, errors.New("invalid or expired token")
	}

	deleteRetweetParam := domain.DeleteRetweetParam{
		UserId:  userSession.UserId,
		TweetId: param.TweetId,
	}

	deleteRetweetResult, errDeleteRetweet := uc.tweetRepository.DeleteRetweet(ctx, deleteRetweetParam)
	if errDeleteRetweet != nil {
		logData["error_delete_retweet"] = errDeleteRetweet.Error()
		uc.logger.
			WithFields(logData).
			WithError(errDeleteRetweet).
			Errorln("error on DeleteRetweet")
		span.End()
		return domain.UndoRetweetUsecaseResult{}, errDeleteRetweet
	}

	logData["delete_retweet_result"] = fmt.Sprintf("%+v", deleteRetweetResult)
	uc.logger.
		WithFields(logData).
		Infoln("success delete retweet")

	span.End()

	result := domain.UndoRetweetUsecaseResult{}
	return result, nil
}
