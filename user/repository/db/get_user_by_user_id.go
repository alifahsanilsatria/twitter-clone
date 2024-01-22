package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *userRepository) GetUserByUserId(ctx context.Context, param domain.GetUserByUserIdParam) (domain.GetUserByUserIdResult, error) {
	logData := logrus.Fields{
		"method":     "userRepository.GetUserByUserId",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	query := `
		select username, email, complete_name, created_at
		from users u
		where u.id = $1
		and u.is_deleted = false
	`

	args := []interface{}{
		param.UserId,
	}

	logData["query"] = query
	logData["args"] = fmt.Sprintf("%+v", args)

	queryRowContextResp := repo.db.QueryRowContext(ctx, query, args...)

	response := domain.GetUserByUserIdResult{}
	errScan := queryRowContextResp.Scan(&response.Username, &response.Email, &response.CompleteName, &response.CreatedAt)
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
		Infoln("success GetUserByUserId")

	return response, nil
}
