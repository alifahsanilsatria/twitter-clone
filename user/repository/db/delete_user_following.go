package db

import (
	"context"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (repo *userRepository) DeleteUserFollowing(ctx context.Context, param domain.DeleteUserFollowingParam) (domain.DeleteUserFollowingResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.DeleteUserFollowing", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

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
		span.End()
		return domain.DeleteUserFollowingResult{}, errScan
	}

	repo.logger.
		WithFields(logData).
		Infoln("success DeleteUserFollowingResult")
	span.End()

	return domain.DeleteUserFollowingResult{}, nil
}
