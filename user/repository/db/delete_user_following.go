package db

import (
	"context"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *userRepository) DeleteUserFollowing(ctx context.Context, param domain.DeleteUserFollowingParam) (domain.DeleteUserFollowingResult, error) {
	logData := logrus.Fields{
		"method":     "userRepository.DeleteUserFollowing",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	query := `
		update user_following
		set is_deleted = true,
		updated_at = $3
		where user_id = $1
		and following_user_id = $2
		returning id
	`

	args := []interface{}{
		param.UserId,
		param.FollowingUserId,
		time.Now(),
	}

	logData["query"] = query
	logData["args"] = fmt.Sprintf("%+v", args)

	queryRowContextResult := repo.db.QueryRowContext(ctx, query, args...)

	result := domain.UpsertUserFollowingResult{}
	errScan := queryRowContextResult.Scan(
		&result.Id,
	)

	if errScan != nil {
		logData["error_scan"] = errScan.Error()
		repo.logger.
			WithFields(logData).
			WithError(errScan).
			Errorln("error on upsert query")
		return domain.DeleteUserFollowingResult{}, errScan
	}

	repo.logger.
		WithFields(logData).
		Infoln("success DeleteUserFollowingResult")

	return domain.DeleteUserFollowingResult{}, nil
}
