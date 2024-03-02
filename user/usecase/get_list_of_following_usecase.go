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

func (uc *userUsecase) GetListOfFollowingUsecase(ctx context.Context, param domain.GetListOfFollowingUsecaseParam) (domain.GetListOfFollowingUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.GetListOfFollowingUsecase", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "userUsecase.GetListOfFollowingUsecase",
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
		return domain.GetListOfFollowingUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.GetListOfFollowingUsecaseResult{}, errors.New("invalid or expired token")
	}

	getListOfFollowingRepoParam := domain.GetListOfFollowingRepoParam{
		UserId: param.UserId,
	}

	logData["get_list_of_following_param"] = fmt.Sprintf("%+v", getListOfFollowingRepoParam)

	getListOfFollowingRepoResult, errGetListOfFollowingRepo := uc.userRepository.GetListOfFollowingRepo(ctx, getListOfFollowingRepoParam)
	if errGetListOfFollowingRepo != nil {
		logData["error_get_list_of_following"] = errGetListOfFollowingRepo.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetListOfFollowingRepo).
			Errorln("error on GetListOfFollowingRepo")
		span.End()
		return domain.GetListOfFollowingUsecaseResult{}, errGetListOfFollowingRepo
	}

	getListOfFollowingUsecaseResult := domain.GetListOfFollowingUsecaseResult{
		Users: make([]domain.GetListOfFollowingUsecaseResult_User, len(getListOfFollowingRepoResult.Users)),
	}

	for idx, data := range getListOfFollowingRepoResult.Users {
		getListOfFollowingUsecaseResult.Users[idx] = domain.GetListOfFollowingUsecaseResult_User{
			UserId:       data.UserId,
			Username:     data.Username,
			CompleteName: data.CompleteName,
		}
	}

	uc.logger.
		WithFields(logData).
		Infoln("success get list of following usecase")
	span.End()

	return getListOfFollowingUsecaseResult, nil
}
