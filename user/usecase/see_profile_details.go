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

func (uc *userUsecase) SeeProfileDetails(ctx context.Context, param domain.SeeProfileDetailsParam) (domain.SeeProfileDetailsResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.SeeProfileDetails", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "userUsecase.SeeProfileDetails",
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
		return domain.SeeProfileDetailsResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.SeeProfileDetailsResult{}, errors.New("invalid or expired token")
	}

	getUserByUserIdParam := domain.GetUserByUserIdParam{
		UserId: userSession.UserId,
	}

	logData["get_user_by_user_id_param"] = fmt.Sprintf("%+v", getUserByUserIdParam)

	getUserByUserIdResult, errGetUserByUserId := uc.userRepository.GetUserByUserId(ctx, getUserByUserIdParam)
	if errGetUserByUserId != nil {
		logData["error_get_user_by_user_id"] = errGetUserByUserId.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserByUserId).
			Errorln("error on GetUserByUserId")
		span.End()
		return domain.SeeProfileDetailsResult{}, errGetUserByUserId
	}

	seeProfileDetailsUsecaseResult := domain.SeeProfileDetailsResult{
		Username:     getUserByUserIdResult.Username,
		Email:        getUserByUserIdResult.CompleteName,
		CompleteName: getUserByUserIdResult.CompleteName,
		CreatedAt:    getUserByUserIdResult.CreatedAt,
	}

	uc.logger.
		WithFields(logData).
		Infoln("success see profile details")
	span.End()

	return seeProfileDetailsUsecaseResult, nil
}
