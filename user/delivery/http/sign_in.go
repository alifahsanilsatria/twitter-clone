package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/labstack/echo"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (handler *userHandler) SignIn(c echo.Context) error {
	logData := logrus.Fields{
		"method": "userHandler.SignIn",
	}

	var reqPayload domain.SignInParam
	errParsingReqPayload := json.NewDecoder(c.Request().Body).Decode(&reqPayload)
	if errParsingReqPayload != nil {
		logData["error_parsing_request_payload"] = errParsingReqPayload.Error()
		handler.logger.
			WithFields(logData).
			WithError(errParsingReqPayload).
			Errorln("error when parsing request payload")
		return c.JSON(http.StatusUnprocessableEntity, errParsingReqPayload.Error())
	}

	logData["request_payload"] = fmt.Sprintf("%+v", reqPayload)

	errvalidateSignInParam := validateSignInParam(reqPayload)
	if errParsingReqPayload != nil {
		logData["error_validate_sign_in_param"] = errvalidateSignInParam.Error()
		handler.logger.
			WithFields(logData).
			WithError(errvalidateSignInParam).
			Errorln("error when validate sign in param")
		return c.JSON(http.StatusBadRequest, errvalidateSignInParam.Error())
	}

	ctx := c.Request().Context()
	signInUsecaseResult, errorSignInUsecase := handler.userUsecase.SignIn(ctx, reqPayload)
	if errorSignInUsecase != nil {
		logData["error_sign_in_usecase"] = errorSignInUsecase.Error()
		handler.logger.
			WithFields(logData).
			WithError(errorSignInUsecase).
			Errorln("error when sign in")
		return c.JSON(http.StatusInternalServerError, errorSignInUsecase.Error())
	}

	return c.JSON(http.StatusOK, signInUsecaseResult)
}

func validateSignInParam(param domain.SignInParam) error {
	if param.Username == "" {
		return errors.New("username is empty")
	}

	if param.Password == "" {
		return errors.New("password is empty")
	}

	return nil
}
