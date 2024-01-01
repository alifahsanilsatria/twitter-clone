package db

import (
	"context"
	"database/sql"
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
		select id, password
		from users u
		where u.username = $1
	`

	args := []interface{}{
		param.Username,
	}

	queryRowContextResp := repo.db.QueryRowContext(ctx, query, args...)

	response := domain.GetUserByUsernameResult{}
	errScan := queryRowContextResp.Scan(&response.Id, &response.HashedPassword)
	if errScan != nil && errScan != sql.ErrNoRows {
		logData["error_scan"] = errScan.Error()
		repo.logger.
			WithFields(logData).
			WithError(errScan).
			Errorln("error on scan")
		return response, errScan
	}

	return response, nil
}
