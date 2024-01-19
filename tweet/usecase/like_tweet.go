package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (uc *tweetUsecase) LikeTweet(ctx context.Context, param domain.LikeTweetUsecaseParam) (domain.LikeTweetUsecaseResult, error) {
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
		return domain.LikeTweetUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
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
		return domain.LikeTweetUsecaseResult{}, errUpsertLikes
	}

	logData["upsert_likes_result"] = fmt.Sprintf("%+v", upsertLikesResult)
	uc.logger.
		WithFields(logData).
		Infoln("success upsert likes")

	likesUsecaseResult := domain.LikeTweetUsecaseResult{}

	return likesUsecaseResult, nil
}
