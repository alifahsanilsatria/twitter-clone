package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (uc *tweetUsecase) PublishTweet(ctx context.Context, param domain.PublishTweetUsecaseParam) (domain.PublishTweetUsecaseResult, error) {
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
		return domain.PublishTweetUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
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
		return domain.PublishTweetUsecaseResult{}, errCreateNewTweet
	}

	logData["create_new_tweet_result"] = fmt.Sprintf("%+v", createNewTweetResult)
	uc.logger.
		WithFields(logData).
		Infoln("success publish tweet")

	publishTweetResult := domain.PublishTweetUsecaseResult{
		TweetId: createNewTweetResult.TweetId,
	}

	return publishTweetResult, nil
}
