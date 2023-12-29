package db

import (
	"context"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *userRepository) CreateNewUserAccount(ctx context.Context, param domain.CreateNewUserAccountParam) (domain.CreateNewUserAccountResult, error) {
	logData := logrus.Fields{
		"method": "userRepository.CreateNewUserAccount",
		"param":  fmt.Sprintf("%+v", param),
	}

	query := `
		insert into user
		(username, password, email, complete_name, is_deleted, created_at)
		values
		($1, $2, $3, $4, $5, $6, $7)
	`

	args := []interface{}{
		param.Username,
		param.HashedPassword,
		param.Email,
		param.CompleteName,
		false,
		time.Now(),
	}

	execContextResp, errQuery := repo.db.ExecContext(ctx, query, args...)
	if errQuery != nil {
		logData["error_query"] = errQuery.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQuery).
			Errorln("error on insert query")
		return domain.CreateNewUserAccountResult{}, errQuery
	}

	idResult, errIdResult := execContextResp.LastInsertId()
	if errIdResult != nil {
		logData["error_id_result"] = errIdResult.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQuery).
			Errorln("error on insert query")
		return domain.CreateNewUserAccountResult{}, errQuery
	}

	response := domain.CreateNewUserAccountResult{
		Id: idResult,
	}

	return response, nil
}
