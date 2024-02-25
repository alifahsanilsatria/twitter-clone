package http

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (handler *userHandler) SignUp(echoCtx echo.Context) error {
	ctx, span := handler.tracer.Start(context.Background(), "handler.SignUp",
		trace.WithSpanKind(trace.SpanKindServer),
	)

	requestId := echoCtx.Request().Header.Get("Request-Id")

	logData := logrus.Fields{
		"method":     "userHandler.SignUp",
		"request_id": requestId,
	}

	ctx = context.WithValue(context.Background(), "request_id", requestId)

	var reqPayload domain.SignUpRequestPayload
	errParsingReqPayload := json.NewDecoder(echoCtx.Request().Body).Decode(&reqPayload)
	if errParsingReqPayload != nil {
		logData["error_parsing_request_payload"] = errParsingReqPayload.Error()
		handler.logger.
			WithFields(logData).
			WithError(errParsingReqPayload).
			Errorln("error when parsing request payload")
		span.End()
		return echoCtx.JSON(http.StatusUnprocessableEntity, errParsingReqPayload.Error())
	}

	span.SetAttributes(attribute.String("request_id", requestId))
	span.SetAttributes(attribute.String("request_payload", fmt.Sprintf("%+v", reqPayload)))

	logData["request_payload"] = fmt.Sprintf("%+v", reqPayload)

	errvalidateSignUpParam := validateSignUpParam(reqPayload)
	if errvalidateSignUpParam != nil {
		logData["error_validate_sign_up_param"] = errvalidateSignUpParam.Error()
		handler.logger.
			WithFields(logData).
			WithError(errvalidateSignUpParam).
			Errorln("error when validate sign up param")
		span.End()
		return echoCtx.JSON(http.StatusBadRequest, errvalidateSignUpParam.Error())
	}

	signUpUsecaseParam := domain.SignUpUsecaseParam{
		Username:     reqPayload.Username,
		Password:     reqPayload.Password,
		Email:        reqPayload.Email,
		CompleteName: reqPayload.CompleteName,
	}

	logData["sign_up_usecase_param"] = fmt.Sprintf("%+v", signUpUsecaseParam)

	signUpUsecaseResult, errorSignUpUsecase := handler.userUsecase.SignUp(ctx, signUpUsecaseParam)
	if errorSignUpUsecase != nil {
		logData["error_sign_up_usecase"] = errorSignUpUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errorSignUpUsecase).
			Errorln("error when parsing request payload")
		span.End()
		return echoCtx.JSON(http.StatusInternalServerError, errorSignUpUsecase.Error())
	}

	logData["sign_up_usecase_result"] = fmt.Sprintf("%+v", signUpUsecaseResult)
	handler.logger.
		WithFields(logData).
		Infoln("success sign up")
	span.End()

	return echoCtx.JSON(http.StatusOK, signUpUsecaseResult)
}

func validateSignUpParam(param domain.SignUpRequestPayload) error {
	if param.Username == "" {
		return errors.New("username is empty")
	}

	if param.Password == "" {
		return errors.New("password is empty")
	}

	if param.Email == "" {
		return errors.New("email is empty")
	}

	if param.CompleteName == "" {
		return errors.New("complete_name is empty")
	}

	usernameRegex, _ := regexp.Compile("^[A-Za-z0-9_]{5,15}$")
	isUsernameMatch := usernameRegex.Match([]byte(param.Username))
	if !isUsernameMatch {
		return errors.New("your username must be more than 4 characters long, can be up to 15 characters, and can only contain letters, numbers, and underscores. no spaces are allowed")
	}

	regexContainUpperCaseLetters, _ := regexp.Compile("[A-Z]")
	regexContainLowerCaseLetters, _ := regexp.Compile("[a-z]")
	regexContainNumbers, _ := regexp.Compile("[0-9]")
	regexContainSymbols, _ := regexp.Compile("[,.\\/!@#$%^&*()_+-=';:]")

	isPasswordMatch := len(param.Password) >= 12 && regexContainUpperCaseLetters.Match([]byte(param.Password)) && regexContainLowerCaseLetters.Match([]byte(param.Password)) && regexContainNumbers.Match([]byte(param.Password)) && regexContainSymbols.Match([]byte(param.Password))
	if !isPasswordMatch {
		return errors.New("your password must be at least 12 characters long and use a mix of uppercase, lowercase, numbers, and symbols character")
	}

	emailRegex, _ := regexp.Compile("^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$")
	isEmailMatch := emailRegex.Match([]byte(param.Email))
	if !isEmailMatch {
		return errors.New("your email is in invalid form")
	}

	if len(param.CompleteName) > 50 {
		return errors.New("your display name must be maximum 50 characters long")
	}

	return nil
}
