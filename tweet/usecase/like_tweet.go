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

func (uc *tweetUsecase) LikeTweet(ctx context.Context, param domain.LikeTweetUsecaseParam) (domain.LikeTweetUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.LikeTweet", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.LikeTweet",
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
		return domain.LikeTweetUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.LikeTweetUsecaseResult{}, errors.New("invalid or expired token")
	}

	upsertLikesParam := domain.UpsertLikesParam{
		UserId:  userSession.UserId,
		TweetId: param.TweetId,
	}

	logData["upsert_likes_param"] = fmt.Sprintf("%+v", upsertLikesParam)

	upsertLikesResult, errUpsertLikes := uc.tweetRepository.UpsertLikes(ctx, upsertLikesParam)
	if errUpsertLikes != nil {
		logData["error_upsert_likes"] = errUpsertLikes.Error()
		uc.logger.
			WithFields(logData).
			WithError(errUpsertLikes).
			Errorln("error on UpsertLikes")
		span.End()
		return domain.LikeTweetUsecaseResult{}, errUpsertLikes
	}

	logData["upsert_likes_result"] = fmt.Sprintf("%+v", upsertLikesResult)
	uc.logger.
		WithFields(logData).
		Infoln("success upsert likes")

	likesUsecaseResult := domain.LikeTweetUsecaseResult{}

	span.End()

	return likesUsecaseResult, nil
}
