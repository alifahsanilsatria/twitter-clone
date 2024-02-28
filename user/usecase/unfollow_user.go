package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (uc *userUsecase) UnfollowUser(ctx context.Context, param domain.UnfollowUserParam) (domain.UnfollowUserResult, error) {
	logData := logrus.Fields{
		"method":     "userUsecase.UnfollowUser",
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
		return domain.UnfollowUserResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		return domain.UnfollowUserResult{}, errors.New("invalid or expired token")
	}

	deleteUserFollowingParam := domain.DeleteUserFollowingParam{
		UserId:          userSession.UserId,
		FollowingUserId: param.FollowingUserId,
	}

	logData["upsert_user_following_param"] = fmt.Sprintf("%+v", deleteUserFollowingParam)

	deleteUserFollowingResult, errDeleteUserFollowing := uc.userRepository.DeleteUserFollowing(ctx, deleteUserFollowingParam)
	if errDeleteUserFollowing != nil {
		logData["error_delete_user_following"] = errDeleteUserFollowing.Error()
		uc.logger.
			WithFields(logData).
			WithError(errDeleteUserFollowing).
			Errorln("error on DeleteUserFollowing")
		return domain.UnfollowUserResult{}, errDeleteUserFollowing
	}

	logData["delete_user_following_result"] = fmt.Sprintf("%+v", deleteUserFollowingResult)

	uc.logger.
		WithFields(logData).
		Infoln("success unfollow user")

	response := domain.UnfollowUserResult{
		Id: deleteUserFollowingResult.Id,
	}

	return response, nil

}
