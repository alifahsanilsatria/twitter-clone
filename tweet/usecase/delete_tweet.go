package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (uc *tweetUsecase) DeleteTweet(ctx context.Context, param domain.DeleteTweetUsecaseParam) (domain.DeleteTweetUsecaseResult, error) {
	logData := logrus.Fields{
		"method":     "tweetUsecase.DeleteTweet",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	getTweetByIdAndUserIdParam := domain.GetTweetByIdAndUserIdParam{
		TweetId: param.TweetId,
		UserId:  param.UserId,
	}

	logData["get_tweet_by_id_and_user_id_param"] = fmt.Sprintf("%+v", getTweetByIdAndUserIdParam)

	getTweetByIdAndUserIdResult, errGetTweetByIdAndUserId := uc.tweetRepository.GetTweetByIdAndUserId(ctx, getTweetByIdAndUserIdParam)
	if errGetTweetByIdAndUserId != nil {
		logData["error_get_tweet_by_id_and_user_id"] = errGetTweetByIdAndUserId.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetTweetByIdAndUserId).
			Errorln("error on GetTweetByIdAndUserId")
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
		return domain.DeleteTweetUsecaseResult{}, errDeleteTweetById
	}

	logData["delete_tweet_by_id_result"] = fmt.Sprintf("%+v", deleteTweetByIdResult)
	uc.logger.
		WithFields(logData).
		Infoln("success delete tweet")

	result := domain.DeleteTweetUsecaseResult{
		TweetId: deleteTweetByIdResult.TweetId,
	}

	return result, nil
}
