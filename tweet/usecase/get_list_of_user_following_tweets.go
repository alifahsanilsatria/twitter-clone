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

func (uc *tweetUsecase) GetListOfUserTimelineTweets(ctx context.Context, param domain.GetListOfUserTimelineTweetsParam) (domain.GetListOfUserTimelineTweetsResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.GetListOfUserTimelineTweets", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.GetListOfUserTimelineTweets",
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
		return domain.GetListOfUserTimelineTweetsResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.GetListOfUserTimelineTweetsResult{}, errors.New("invalid or expired token")
	}

	getListOfUserFollowingTweetsParam := domain.GetListOfUserFollowingTweetsParam{
		UserId: userSession.UserId,
	}
	getListOfUserFollowingTweetsResult, errGetListOfUserFollowingTweets := uc.tweetRepository.GetListOfUserFollowingTweets(ctx, getListOfUserFollowingTweetsParam)
	if errGetListOfUserFollowingTweets != nil {
		logData["error_get_list_of_user_following_tweets"] = errGetListOfUserFollowingTweets.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetListOfUserFollowingTweets).
			Errorln("error on GetListOfUserFollowingTweets")
		span.End()
		return domain.GetListOfUserTimelineTweetsResult{}, errGetListOfUserFollowingTweets
	}

	ucResult := domain.GetListOfUserTimelineTweetsResult{
		Tweets: make([]domain.GetListOfUserTimelineTweetsResult_TweetDetails, len(getListOfUserFollowingTweetsResult.Tweets)),
	}

	for idx, tweet := range getListOfUserFollowingTweetsResult.Tweets {
		ucResult.Tweets[idx] = domain.GetListOfUserTimelineTweetsResult_TweetDetails{
			TweetId:      tweet.TweetId,
			Username:     tweet.Username,
			CompleteName: tweet.CompleteName,
			Content:      tweet.Content,
			CountRetweet: tweet.CountRetweet,
			CountLikes:   tweet.CountLikes,
			CountReplies: tweet.CountReplies,
		}
	}

	span.End()

	return ucResult, nil
}
