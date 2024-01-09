package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (uc *tweetUsecase) DeleteTweet(ctx context.Context, param domain.DeleteTweetParam) (domain.DeleteTweetResult, error) {
	logData := logrus.Fields{
		"method": "tweetUsecase.DeleteTweet",
		"param":  fmt.Sprintf("%+v", param),
	}

	getTweetByIdAndUserIdParam := domain.GetTweetByIdAndUserIdParam{
		TweetId: param.TweetId,
		UserId:  param.UserId,
	}
	getTweetByIdAndUserIdResult, errGetTweetByIdAndUserId := uc.tweetRepository.GetTweetByIdAndUserId(ctx, getTweetByIdAndUserIdParam)
	if errGetTweetByIdAndUserId != nil {
		logData["error_get_tweet_by_id_and_user_id"] = errGetTweetByIdAndUserId.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetTweetByIdAndUserId).
			Errorln("error on GetTweetByIdAndUserId")
		return domain.DeleteTweetResult{}, errGetTweetByIdAndUserId
	}

	if getTweetByIdAndUserIdResult.TweetId == 0 {
		return domain.DeleteTweetResult{}, errors.New("invalid user or deleted tweet")
	}

	deleteTweetByIdParam := domain.DeleteTweetByIdParam{
		TweetId: param.TweetId,
	}
	deleteTweetById, errDeleteTweetById := uc.tweetRepository.DeleteTweetById(ctx, deleteTweetByIdParam)
	if errDeleteTweetById != nil {
		logData["error_delete_tweet_by_id"] = errDeleteTweetById.Error()
		uc.logger.
			WithFields(logData).
			WithError(errDeleteTweetById).
			Errorln("error on DeleteTweetById")
		return domain.DeleteTweetResult{}, errDeleteTweetById
	}

	result := domain.DeleteTweetResult{
		TweetId: deleteTweetById.TweetId,
	}

	return result, nil
}
