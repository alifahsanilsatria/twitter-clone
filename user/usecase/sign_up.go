package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/alifahsanilsatria/twitter-clone/common"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (uc *userUsecase) SignUp(ctx context.Context, param domain.SignUpUsecaseParam) (domain.SignUpResult, error) {
	ctx, span := uc.tracer.Start(ctx, "usecase.SignUp", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "userUsecase.SignUp",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	getUserByUsernameOrEmailParam := domain.GetUserByUsernameOrEmailParam{
		Username: strings.ToLower(param.Username),
		Email:    param.Email,
	}

	logData["get_user_by_username_or_email_param"] = fmt.Sprintf("%+v", getUserByUsernameOrEmailParam)

	getUserByUsernameOrEmailResp, errGetUserByUsernameOrEmail := uc.userRepository.GetUserByUsernameOrEmail(ctx, getUserByUsernameOrEmailParam)
	if errGetUserByUsernameOrEmail != nil {
		logData["error_get_user_by_username_or_email"] = errGetUserByUsernameOrEmail.Error()
		uc.logger.
			WithFields(logData).
			WithError(errGetUserByUsernameOrEmail).
			Errorln("error on GetUserByUsernameOrEmail")
		span.End()
		return domain.SignUpResult{}, errGetUserByUsernameOrEmail
	}

	logData["get_user_by_username_or_email_result"] = fmt.Sprintf("%+v", getUserByUsernameOrEmailResp)

	if getUserByUsernameOrEmailResp.Id > 0 {
		err := errors.New("username or email already exists")
		span.End()
		return domain.SignUpResult{}, err
	}

	bytesHashedPassword, errHashPassword := common.HashPlaintext(param.Password)
	if errHashPassword != nil {
		logData["error_hash_password"] = errHashPassword.Error()
		uc.logger.
			WithFields(logData).
			WithError(errHashPassword).
			Errorln("error on HashPassword")
		span.End()
		return domain.SignUpResult{}, errHashPassword
	}

	hashedPassword := string(bytesHashedPassword)
	createNewUserAccountParam := domain.CreateNewUserAccountParam{
		Username:       param.Username,
		HashedPassword: hashedPassword,
		Email:          param.Email,
		CompleteName:   param.CompleteName,
	}

	logData["create_new_user_account_param"] = fmt.Sprintf("%+v", createNewUserAccountParam)

	createNewUserAccountResp, errCreateNewUserAccount := uc.userRepository.CreateNewUserAccount(ctx, createNewUserAccountParam)
	if errCreateNewUserAccount != nil {
		logData["error_create_new_user_account"] = errCreateNewUserAccount.Error()
		uc.logger.
			WithFields(logData).
			WithError(errCreateNewUserAccount).
			Errorln("error on CreateNewUserAccount")
		span.End()
		return domain.SignUpResult{}, errCreateNewUserAccount
	}

	logData["create_new_user_account_resp"] = fmt.Sprintf("%+v", createNewUserAccountResp)

	uc.logger.
		WithFields(logData).
		Infoln("success sign up")

	response := domain.SignUpResult{
		Id: createNewUserAccountResp.Id,
	}

	span.End()

	return response, nil

}
