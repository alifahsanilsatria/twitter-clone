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

func (repo *userRepository) UpsertUserFollowing(ctx context.Context, param domain.UpsertUserFollowingParam) (domain.UpsertUserFollowingResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.UpsertUserFollowing", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "userRepository.UpsertUserFollowing",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	query := `
		insert into user_following
		(user_id, following_user_id, is_deleted, created_at)
		values
		($1, $2, false, $3)
		on conflict (user_id, following_user_id)
		do update
		set is_deleted = false,
		updated_at = $3
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
		return domain.UpsertUserFollowingResult{}, errScan
	}

	repo.logger.
		WithFields(logData).
		Infoln("success UpsertUserFollowing")
	span.End()

	return domain.UpsertUserFollowingResult{}, nil
}
