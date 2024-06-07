package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (uc *tweetUsecase) GetListOfAUserTimelineTweets(ctx context.Context, param domain.GetListOfAUserTimelineTweetsParam) (domain.GetListOfAUserTimelineTweetsResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.GetListOfUserTimelineTweets", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetUsecase.GetListOfAUserTimelineTweets",
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
		return domain.GetListOfAUserTimelineTweetsResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.GetListOfAUserTimelineTweetsResult{}, errors.New("invalid or expired token")
	}

	getUserByUsernameParam := domain.GetUserByUsernameParam{
		Username: strings.ToLower(param.Username),
	}
	user, errGetUserByUsername := uc.userRepository.GetUserByUsername(ctx, getUserByUsernameParam)
	if errGetUserByUsername != nil {
		logData["error_get_user_by_username"] = errGetUserByUsername.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserByUsername).
			Errorln("error on GetUserByUsername")
		span.End()
		return domain.GetListOfAUserTimelineTweetsResult{}, errGetUserByUsername
	}

	getListOfAUserTweetsParam := domain.GetListOfAUserTweetsParam{
		UserId: user.Id,
	}
	getListOfAUserTweetsResult, errGetListOfAUserTweets := uc.tweetRepository.GetListOfAUserTweets(ctx, getListOfAUserTweetsParam)
	if errGetListOfAUserTweets != nil {
		logData["error_get_list_of_a_user_tweets"] = errGetListOfAUserTweets.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetListOfAUserTweets).
			Errorln("error on GetListOfAUserTweets")
		span.End()
		return domain.GetListOfAUserTimelineTweetsResult{}, errGetListOfAUserTweets
	}

	ucResult := domain.GetListOfAUserTimelineTweetsResult{
		Tweets: make([]domain.GetListOfAUserTimelineTweetsResult_TweetDetails, len(getListOfAUserTweetsResult.Tweets)),
	}

	for idx, tweet := range getListOfAUserTweetsResult.Tweets {
		ucResult.Tweets[idx] = domain.GetListOfAUserTimelineTweetsResult_TweetDetails{
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
