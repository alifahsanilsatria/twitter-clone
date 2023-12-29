package db

import (
	"context"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *userRepository) GetUserByUsernameOrEmail(ctx context.Context, param domain.GetUserByUsernameOrEmailParam) (domain.GetUserByUsernameOrEmailResult, error) {
	logData := logrus.Fields{
		"method": "userRepository.GetUserByUsernameOrEmail",
		"param":  fmt.Sprintf("%+v", param),
	}
	query := `
		select id
		from user u
		where username = $1
		or email = $2
	`

	args := []interface{}{
		param.Username,
		param.Email,
	}

	queryRowContextResp := repo.db.QueryRowContext(ctx, query, args...)

	response := domain.GetUserByUsernameOrEmailResult{}
	errScan := queryRowContextResp.Scan(&response.Id)
	if errScan != nil {
		logData["error_scan"] = errScan.Error()
		repo.logger.
			WithFields(logData).
			WithError(errScan).
			Errorln("error on scan")
	} else {
		logData["response"] = fmt.Sprintf("%+v", response)
		repo.logger.
			WithFields(logData).
			WithError(errScan).
			Debugln("success get response")
	}

	return response, errScan
}
