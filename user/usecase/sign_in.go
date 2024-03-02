package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/crypto/bcrypt"
)

func (uc *userUsecase) SignIn(ctx context.Context, param domain.SignInUsecaseParam) (domain.SignInResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.SignIn", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "userUsecase.SignIn",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	getUserByUsernameParam := domain.GetUserByUsernameParam{
		Username: param.Username,
	}

	logData["get_user_by_username_param"] = fmt.Sprintf("%+v", getUserByUsernameParam)

	getUserByUsernameResult, errGetUserByUsername := uc.userRepository.GetUserByUsername(ctx, getUserByUsernameParam)
	if errGetUserByUsername != nil {
		logData["error_get_user_by_username"] = errGetUserByUsername.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserByUsername).
			Errorln("error on GetUserByUsername")
		span.End()
		return domain.SignInResult{}, errGetUserByUsername
	}

	logData["get_user_by_username_result"] = fmt.Sprintf("%+v", getUserByUsernameResult)

	errPasswordValidation := bcrypt.CompareHashAndPassword([]byte(getUserByUsernameResult.HashedPassword), []byte(param.Password))
	if errPasswordValidation != nil {
		logData["error_compare_hash_and_password"] = errPasswordValidation.Error()
		uc.logger.
			WithFields(logData).
			WithError(errPasswordValidation).
			Errorln("error on CompareHashAndPassword")
		span.End()
		return domain.SignInResult{}, errors.New("your username or password is wrong")
	}

	token, _ := uuid.NewRandom()
	setUserSessionToRedisParam := domain.SetUserSessionToRedisParam{
		UserId: getUserByUsernameResult.Id,
		Token:  token.String(),
		TTL:    30 * 24 * time.Hour,
	}

	logData["set_user_session_to_redis_param"] = fmt.Sprintf("%+v", setUserSessionToRedisParam)

	errUpsertUserSession := uc.userSessionRepository.SetUserSessionToRedis(ctx, setUserSessionToRedisParam)
	if errUpsertUserSession != nil {
		logData["error_upsert_user_session"] = errUpsertUserSession.Error()
		uc.logger.
			WithFields(logData).
			WithError(errUpsertUserSession).
			Errorln("error on UpsertUserSession")
		span.End()
		return domain.SignInResult{}, errUpsertUserSession
	}

	uc.logger.
		WithFields(logData).
		Infoln("success sign in")

	result := domain.SignInResult{
		Token: token.String(),
	}

	span.End()
	return result, nil

}
