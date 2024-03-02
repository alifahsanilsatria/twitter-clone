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

func (uc *tweetUsecase) DeleteTweet(ctx context.Context, param domain.DeleteTweetUsecaseParam) (domain.DeleteTweetUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.DeleteTweet", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.DeleteTweet",
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
		return domain.DeleteTweetUsecaseResult{}, errGetUserSession
	}

	if userSession.UserId == 0 {
		span.End()
		return domain.DeleteTweetUsecaseResult{}, errors.New("invalid or expired token")
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		return domain.DeleteTweetUsecaseResult{}, errors.New("invalid or expired token")
	}

	getTweetByIdAndUserIdParam := domain.GetTweetByIdAndUserIdParam{
		TweetId: param.TweetId,
		UserId:  userSession.UserId,
	}

	logData["get_tweet_by_id_and_user_id_param"] = fmt.Sprintf("%+v", getTweetByIdAndUserIdParam)

	getTweetByIdAndUserIdResult, errGetTweetByIdAndUserId := uc.tweetRepository.GetTweetByIdAndUserId(ctx, getTweetByIdAndUserIdParam)
	if errGetTweetByIdAndUserId != nil {
		logData["error_get_tweet_by_id_and_user_id"] = errGetTweetByIdAndUserId.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetTweetByIdAndUserId).
			Errorln("error on GetTweetByIdAndUserId")
		span.End()
		return domain.DeleteTweetUsecaseResult{}, errGetTweetByIdAndUserId
	}

	logData["get_tweet_by_id_and_user_id_result"] = fmt.Sprintf("%+v", getTweetByIdAndUserIdResult)

	if getTweetByIdAndUserIdResult.TweetId == 0 {
		return domain.DeleteTweetUsecaseResult{}, errors.New("invalid user or deleted tweet")
	}

	deleteTweetByIdParam := domain.DeleteTweetByIdParam{
		TweetId: param.TweetId,
	}

	logData["delete_tweet_by_id_param"] = fmt.Sprintf("%+v", deleteTweetByIdParam)

	deleteTweetByIdResult, errDeleteTweetById := uc.tweetRepository.DeleteTweetById(ctx, deleteTweetByIdParam)
	if errDeleteTweetById != nil {
		logData["error_delete_tweet_by_id"] = errDeleteTweetById.Error()
		uc.logger.
			WithFields(logData).
			WithError(errDeleteTweetById).
			Errorln("error on DeleteTweetById")
		span.End()
		return domain.DeleteTweetUsecaseResult{}, errDeleteTweetById
	}

	logData["delete_tweet_by_id_result"] = fmt.Sprintf("%+v", deleteTweetByIdResult)
	uc.logger.
		WithFields(logData).
		Infoln("success delete tweet")

	result := domain.DeleteTweetUsecaseResult{
		TweetId: deleteTweetByIdResult.TweetId,
	}

	span.End()

	return result, nil
}
