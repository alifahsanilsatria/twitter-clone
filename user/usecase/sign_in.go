package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func (uc *userUsecase) SignIn(ctx context.Context, param domain.SignInParam) (domain.SignInResult, error) {
	logData := logrus.Fields{
		"method": "userUsecase.SignIn",
		"param":  fmt.Sprintf("%+v", param),
	}

	getUserByUsernameParam := domain.GetUserByUsernameParam{}
	getUserByUsernameResult, errGetUserByUsername := uc.userRepository.GetUserByUsername(ctx, getUserByUsernameParam)
	if errGetUserByUsername != nil {
		logData["error_get_user_by_username"] = errGetUserByUsername.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserByUsername).
			Errorln("error on GetUserByUsername")
		return domain.SignInResult{}, errGetUserByUsername
	}

	errPasswordValidation := bcrypt.CompareHashAndPassword([]byte(getUserByUsernameResult.HashedPassword), []byte(param.Password))
	if errPasswordValidation != nil {
		logData["error_compare_hash_and_password"] = errPasswordValidation.Error()
		uc.logger.
			WithFields(logData).
			WithError(errPasswordValidation).
			Errorln("error on CompareHashAndPassword")
		return domain.SignInResult{}, errors.New("your username or password is wrong")
	}

	token, _ := uuid.NewRandom()
	setUserSessionToRedisParam := domain.SetUserSessionToRedisParam{
		UserId: getUserByUsernameResult.Id,
		Token:  token.String(),
		TTL:    30 * 24 * time.Hour,
	}
	errUpsertUserSession := uc.userSessionRepository.SetUserSessionToRedis(ctx, setUserSessionToRedisParam)
	if errUpsertUserSession != nil {
		logData["error_upsert_user_session"] = errUpsertUserSession.Error()
		uc.logger.
			WithFields(logData).
			WithError(errUpsertUserSession).
			Errorln("error on UpsertUserSession")
		return domain.SignInResult{}, errUpsertUserSession
	}

	result := domain.SignInResult{
		Token: token.String(),
	}

	return result, nil

}
