package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (uc *userUsecase) SignOut(ctx context.Context, param domain.SignOutParam) (domain.SignOutResult, error) {
	logData := logrus.Fields{
		"method": "userUsecase.SignOut",
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
		return domain.SignOutResult{}, errGetUserSession
	}

	if userSession.UserId == 0 {
		return domain.SignOutResult{}, errors.New("invalid or expired token")
	}

	deleteSessionByTokenParam := domain.DeleteUserSessionByTokenParam{
		Token: param.Token,
	}
	errDeleteSession := uc.userSessionRepository.DeleteUserSessionByToken(ctx, deleteSessionByTokenParam)
	if errDeleteSession != nil {
		logData["error_delete_user_session"] = errDeleteSession.Error()
		uc.logger.
			WithFields(logData).
			WithError(errDeleteSession).
			Errorln("error on DeleteUserSessionByToken")
		return domain.SignOutResult{}, errDeleteSession
	}

	return domain.SignOutResult{}, nil
}
