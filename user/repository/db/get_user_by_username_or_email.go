package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *userRepository) GetUserByUsernameOrEmail(ctx context.Context, param domain.GetUserByUsernameOrEmailParam) (domain.GetUserByUsernameOrEmailResult, error) {
	logData := logrus.Fields{
		"method":     "userRepository.GetUserByUsernameOrEmail",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}
	query := `
		select id
		from users u
		where u.username = $1
		or u.email = $2
	`

	args := []interface{}{
		param.Username,
		param.Email,
	}

	logData["query"] = query
	logData["args"] = fmt.Sprintf("%+v", args)

	queryRowContextResp := repo.db.QueryRowContext(ctx, query, args...)

	response := domain.GetUserByUsernameOrEmailResult{}
	errScan := queryRowContextResp.Scan(&response.Id)
	if errScan != nil && errScan != sql.ErrNoRows {
		logData["error_scan"] = errScan.Error()
		repo.logger.
			WithFields(logData).
			WithError(errScan).
			Errorln("error on scan")
		return response, errScan
	}

	logData["response"] = fmt.Sprintf("%+v", response)
	repo.logger.
		WithFields(logData).
		Infoln("success GetUserByUsernameOrEmail")

	return response, nil
}
