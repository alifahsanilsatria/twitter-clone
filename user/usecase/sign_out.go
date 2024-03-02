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

func (uc *userUsecase) SignOut(ctx context.Context, param domain.SignOutUsecaseParam) (domain.SignOutResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.SignOut", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "userUsecase.SignOut",
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
		return domain.SignOutResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.SignOutResult{}, errors.New("invalid or expired token")
	}

	deleteSessionByTokenParam := domain.DeleteUserSessionByTokenParam{
		Token: param.Token,
	}

	logData["delete_session_by_token_param"] = fmt.Sprintf("%+v", deleteSessionByTokenParam)

	errDeleteSession := uc.userSessionRepository.DeleteUserSessionByToken(ctx, deleteSessionByTokenParam)
	if errDeleteSession != nil {
		logData["error_delete_user_session"] = errDeleteSession.Error()
		uc.logger.
			WithFields(logData).
			WithError(errDeleteSession).
			Errorln("error on DeleteUserSessionByToken")
		span.End()
		return domain.SignOutResult{}, errDeleteSession
	}

	uc.logger.
		WithFields(logData).
		Infoln("success sign out")
	span.End()

	return domain.SignOutResult{}, nil
}
