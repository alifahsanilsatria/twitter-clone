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

func (uc *tweetUsecase) PublishTweet(ctx context.Context, param domain.PublishTweetUsecaseParam) (domain.PublishTweetUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.PublishTweet", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.PublishTweet",
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
		return domain.PublishTweetUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.PublishTweetUsecaseResult{}, errors.New("invalid or expired token")
	}

	createNewTweetParam := domain.CreateNewTweetParam{
		UserId:   userSession.UserId,
		ParentId: param.ParentId,
		Content:  param.Content,
	}

	logData["create_new_tweet_param"] = fmt.Sprintf("%+v", createNewTweetParam)

	createNewTweetResult, errCreateNewTweet := uc.tweetRepository.CreateNewTweet(ctx, createNewTweetParam)
	if errCreateNewTweet != nil {
		logData["error_create_new_tweet"] = errCreateNewTweet.Error()
		uc.logger.
			WithFields(logData).
			WithError(errCreateNewTweet).
			Errorln("error on CreateNewTweet")
		span.End()
		return domain.PublishTweetUsecaseResult{}, errCreateNewTweet
	}

	logData["create_new_tweet_result"] = fmt.Sprintf("%+v", createNewTweetResult)
	uc.logger.
		WithFields(logData).
		Infoln("success publish tweet")

	publishTweetResult := domain.PublishTweetUsecaseResult{
		TweetId: createNewTweetResult.TweetId,
	}

	span.End()
	return publishTweetResult, nil
}
