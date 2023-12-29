package db

import (
	"context"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *userRepository) GetUserByUsername(ctx context.Context, param domain.GetUserByUsernameParam) (domain.GetUserByUsernameResult, error) {
	logData := logrus.Fields{
		"method": "userRepository.GetUserByUsername",
		"param":  fmt.Sprintf("%+v", param),
	}

	query := `
		select id
		from user u
		where username = $1
	`

	args := []interface{}{
		param.Username,
	}

	queryRowContextResp := repo.db.QueryRowContext(ctx, query, args...)

	response := domain.GetUserByUsernameResult{}
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
