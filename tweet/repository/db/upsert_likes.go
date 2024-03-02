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

func (repo *tweetRepository) UpsertLikes(ctx context.Context, param domain.UpsertLikesParam) (domain.UpsertLikesResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.UpsertLikes", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetRepository.UpsertLikes",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	queryUpsertLikes := `
		insert into likes
		(user_id, tweet_id, is_deleted, created_at)
		values
		($1, $2, false, $3)
		on conflict (user_id, tweet_id)
		do update
		set is_deleted = false,
		updated_at = $3
	`

	argsQueryUpsertLikes := []interface{}{
		param.UserId,
		param.TweetId,
		time.Now(),
	}

	logData["query_upsert_likes"] = queryUpsertLikes
	logData["args_query_upsert_likes"] = fmt.Sprintf("%+v", argsQueryUpsertLikes)

	errUpsertQuery := repo.db.QueryRowContext(ctx, queryUpsertLikes, argsQueryUpsertLikes...).Err()
	if errUpsertQuery != nil {
		logData["error_query_upsert_likes"] = errUpsertQuery.Error()
		repo.logger.
			WithFields(logData).
			WithError(errUpsertQuery).
			Errorln("error on upsert query")
		span.End()
		return domain.UpsertLikesResult{}, errUpsertQuery
	}

	span.End()
	return domain.UpsertLikesResult{}, nil
}
