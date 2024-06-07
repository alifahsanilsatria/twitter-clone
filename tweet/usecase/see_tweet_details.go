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

func (uc *tweetUsecase) SeeTweetDetails(ctx context.Context, param domain.SeeTweetDetailsUsecaseParam) (domain.SeeTweetDetailsUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.GetListOfUserTimelineTweets", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.SeeTweetDetails",
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
		return domain.SeeTweetDetailsUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.SeeTweetDetailsUsecaseResult{}, errors.New("invalid or expired token")
	}

	getTweetByIdParam := domain.GetTweetByIdParam{
		TweetId: param.TweetId,
	}
	getTweetByIdResult, errGetTweetById := uc.tweetRepository.GetTweetById(ctx, getTweetByIdParam)
	if errGetTweetById != nil {
		logData["error_get_tweet_by_id"] = errGetTweetById.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetTweetById).
			Errorln("error on GetTweetById")
		span.End()
		return domain.SeeTweetDetailsUsecaseResult{}, errGetTweetById
	}

	if getTweetByIdResult.TweetId == 0 {
		span.End()
		return domain.SeeTweetDetailsUsecaseResult{}, errors.New("the tweet doesn't exist or deleted")
	}

	getParentsDataByTweetIdParam := domain.GetParentsDataByTweetIdParam{
		TweetId: param.TweetId,
	}
	getParentsDataByTweetIdResult, errGetParentsDataByTweetId := uc.tweetRepository.GetParentsDataByTweetId(ctx, getParentsDataByTweetIdParam)
	if errGetParentsDataByTweetId != nil {
		logData["error_get_parents_data_by_tweet_id"] = errGetParentsDataByTweetId.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetParentsDataByTweetId).
			Errorln("error on GetParentsDataByTweetId")
		span.End()
		return domain.SeeTweetDetailsUsecaseResult{}, errGetParentsDataByTweetId
	}

	getChildrenDataByTweetIdParam := domain.GetChildrenDataByTweetIdParam{
		TweetId: param.TweetId,
	}
	getChildrenDataByTweetIdResult, errGetChildrenDataByTweetId := uc.tweetRepository.GetChildrenDataByTweetId(ctx, getChildrenDataByTweetIdParam)
	if errGetChildrenDataByTweetId != nil {
		logData["error_get_children_data_by_tweet_id"] = errGetChildrenDataByTweetId.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetChildrenDataByTweetId).
			Errorln("error on GetChildrenDataByTweetId")
		span.End()
		return domain.SeeTweetDetailsUsecaseResult{}, errGetChildrenDataByTweetId
	}

	seeTweetDetailsUsecaseResult := domain.SeeTweetDetailsUsecaseResult{}

	for _, parent := range getParentsDataByTweetIdResult.Parent {
		tweetDetails := domain.SeeTweetDetailsUsecaseResult_TweetDetails{
			TweetId:       parent.TweetId,
			Username:      parent.Username,
			CompleteName:  parent.CompleteName,
			Content:       parent.Content,
			CountRetweet:  parent.CountRetweet,
			CountLikes:    parent.CountLikes,
			CountReplies:  parent.CountReplies,
			IsParentTweet: true,
		}
		seeTweetDetailsUsecaseResult.TweetDetails = append(seeTweetDetailsUsecaseResult.TweetDetails, tweetDetails)
	}

	tweetDetails := domain.SeeTweetDetailsUsecaseResult_TweetDetails{
		TweetId:        getTweetByIdResult.TweetId,
		Username:       getTweetByIdResult.Username,
		CompleteName:   getTweetByIdResult.CompleteName,
		Content:        getTweetByIdResult.Content,
		CountRetweet:   getTweetByIdResult.CountRetweet,
		CountLikes:     getTweetByIdResult.CountLikes,
		CountReplies:   getTweetByIdResult.CountReplies,
		IsCurrentTweet: true,
	}
	seeTweetDetailsUsecaseResult.TweetDetails = append(seeTweetDetailsUsecaseResult.TweetDetails, tweetDetails)

	for _, child := range getChildrenDataByTweetIdResult.ChildTweet {
		tweetDetails := domain.SeeTweetDetailsUsecaseResult_TweetDetails{
			TweetId:      child.TweetId,
			Username:     child.Username,
			CompleteName: child.CompleteName,
			Content:      child.Content,
			CountRetweet: child.CountRetweet,
			CountLikes:   child.CountLikes,
			CountReplies: child.CountReplies,
			IsChildTweet: true,
		}
		seeTweetDetailsUsecaseResult.TweetDetails = append(seeTweetDetailsUsecaseResult.TweetDetails, tweetDetails)
	}

	span.End()

	return seeTweetDetailsUsecaseResult, nil
}
