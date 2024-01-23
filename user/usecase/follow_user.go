package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (uc *userUsecase) FollowUser(ctx context.Context, param domain.FollowUserParam) (domain.FollowUserResult, error) {
	logData := logrus.Fields{
		"method":     "userUsecase.FollowUser",
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
		return domain.FollowUserResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		return domain.FollowUserResult{}, errors.New("invalid or expired token")
	}

	getUserByUserIdParam := domain.GetUserByUserIdParam{
		UserId: param.FollowingUserId,
	}

	logData["get_user_by_user_id_param"] = fmt.Sprintf("%+v", getUserByUserIdParam)

	getUserByUserIdResult, errGetUserByUserId := uc.userRepository.GetUserByUserId(ctx, getUserByUserIdParam)
	if errGetUserByUserId != nil {
		logData["error_get_user_by_user_id"] = errGetUserByUserId.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserByUserId).
			Errorln("error on GetUserByUserId")
		return domain.FollowUserResult{}, errGetUserByUserId
	}

	upsertUserFollowingParam := domain.UpsertUserFollowingParam{
		UserId:          userSession.UserId,
		FollowingUserId: getUserByUserIdResult.UserId,
	}

	logData["upsert_user_following_param"] = fmt.Sprintf("%+v", upsertUserFollowingParam)

	upsertUserFollowingResult, errUpsertUserFollowing := uc.userRepository.UpsertUserFollowing(ctx, upsertUserFollowingParam)
	if errUpsertUserFollowing != nil {
		logData["error_upsert_user_following"] = errUpsertUserFollowing.Error()
		uc.logger.
			WithFields(logData).
			WithError(errUpsertUserFollowing).
			Errorln("error on UpsertUserFollowing")
		return domain.FollowUserResult{}, errUpsertUserFollowing
	}

	logData["upsert_user_following_result"] = fmt.Sprintf("%+v", upsertUserFollowingResult)

	uc.logger.
		WithFields(logData).
		Infoln("success follow user")

	response := domain.FollowUserResult{
		Id: upsertUserFollowingResult.Id,
	}

	return response, nil

}
