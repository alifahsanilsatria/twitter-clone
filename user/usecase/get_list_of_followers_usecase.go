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

func (uc *userUsecase) GetListOfFollowersUsecase(ctx context.Context, param domain.GetListOfFollowersUsecaseParam) (domain.GetListOfFollowersUsecaseResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.GetListOfFollowersUsecase", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "userUsecase.GetListOfFollowersUsecase",
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
		return domain.GetListOfFollowersUsecaseResult{}, errGetUserSession
	}

	logData["get_user_session_by_token_result"] = fmt.Sprintf("%+v", userSession)

	if userSession.UserId == 0 {
		span.End()
		return domain.GetListOfFollowersUsecaseResult{}, errors.New("invalid or expired token")
	}

	getListOfFollowersRepoParam := domain.GetListOfFollowersRepoParam{
		UserId: param.UserId,
	}

	logData["get_list_of_followers_param"] = fmt.Sprintf("%+v", getListOfFollowersRepoParam)

	getListOfFollowersRepoResult, errGetListOfFollowersRepo := uc.userRepository.GetListOfFollowersRepo(ctx, getListOfFollowersRepoParam)
	if errGetListOfFollowersRepo != nil {
		logData["error_get_list_of_followers"] = errGetListOfFollowersRepo.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetListOfFollowersRepo).
			Errorln("error on GetListOfFollowersRepo")
		span.End()
		return domain.GetListOfFollowersUsecaseResult{}, errGetListOfFollowersRepo
	}

	getListOfFollowersUsecaseResult := domain.GetListOfFollowersUsecaseResult{
		Users: make([]domain.GetListOfFollowersUsecaseResult_User, len(getListOfFollowersRepoResult.Users)),
	}

	for idx, data := range getListOfFollowersRepoResult.Users {
		getListOfFollowersUsecaseResult.Users[idx] = domain.GetListOfFollowersUsecaseResult_User{
			UserId:       data.UserId,
			Username:     data.Username,
			CompleteName: data.CompleteName,
		}
	}

	uc.logger.
		WithFields(logData).
		Infoln("success get list of followers usecase")
	span.End()

	return getListOfFollowersUsecaseResult, nil
}
