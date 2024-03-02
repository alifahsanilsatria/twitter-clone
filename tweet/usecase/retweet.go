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

func (uc *tweetUsecase) Retweet(ctx context.Context, param domain.RetweetUsecaseParam) (domain.RetweetUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.Retweet", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.Retweet",
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
		return domain.RetweetUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.RetweetUsecaseResult{}, errors.New("invalid or expired token")
	}

	upsertRetweetParam := domain.UpsertRetweetParam{
		UserId:  userSession.UserId,
		TweetId: param.TweetId,
	}

	logData["upsert_retweet_param"] = fmt.Sprintf("%+v", upsertRetweetParam)

	upsertRetweetResult, errUpsertRetweet := uc.tweetRepository.UpsertRetweet(ctx, upsertRetweetParam)
	if errUpsertRetweet != nil {
		logData["error_upsert_retweet"] = errUpsertRetweet.Error()
		uc.logger.
			WithFields(logData).
			WithError(errUpsertRetweet).
			Errorln("error on UpsertRetweet")
		span.End()
		return domain.RetweetUsecaseResult{}, errUpsertRetweet
	}

	logData["upsert_retweet_result"] = fmt.Sprintf("%+v", upsertRetweetResult)
	uc.logger.
		WithFields(logData).
		Infoln("success upsert retweet")

	retweetUsecaseResult := domain.RetweetUsecaseResult{}

	span.End()

	return retweetUsecaseResult, nil
}
