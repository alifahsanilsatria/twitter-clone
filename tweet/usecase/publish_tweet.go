package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (uc *tweetUsecase) PublishTweet(ctx context.Context, param domain.PublishTweetParam) (domain.PublishTweetResult, error) {
	logData := logrus.Fields{
		"method": "tweetUsecase.PublishTweet",
		"param":  fmt.Sprintf("%+v", param),
	}

	getUserSessionByTokenParam := domain.GetUserSessionByTokenParam{
		Token: param.Token,
	}
	userSession, errGetUserSession := uc.userSessionRepository.GetUserSessionByToken(ctx, getUserSessionByTokenParam)
	if errGetUserSession != nil {
		logData["error_get_user_session"] = errGetUserSession.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserSession).
			Errorln("error on GetUserSessionByToken")
		return domain.PublishTweetResult{}, errGetUserSession
	}

	if userSession.UserId == 0 {
		return domain.PublishTweetResult{}, errors.New("invalid or expired token")
	}

	createNewTweetParam := domain.CreateNewTweetParam{
		UserId:   userSession.UserId,
		ParentId: param.ParentId,
		Content:  param.Content,
	}
	createNewTweetResult, errCreateNewTweet := uc.tweetRepository.CreateNewTweet(ctx, createNewTweetParam)
	if errCreateNewTweet != nil {
		logData["error_create_new_tweet"] = errCreateNewTweet.Error()
		uc.logger.
			WithFields(logData).
			WithError(errCreateNewTweet).
			Errorln("error on CreateNewTweet")
		return domain.PublishTweetResult{}, errCreateNewTweet
	}

	publishTweetResult := domain.PublishTweetResult{
		TweetId: createNewTweetResult.TweetId,
	}

	return publishTweetResult, nil
}
