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

func (repo *tweetRepository) UpsertRetweet(ctx context.Context, param domain.UpsertRetweetParam) (domain.UpsertRetweetResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.UpsertRetweet", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetRepository.UpsertRetweet",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	queryUpsertRetweet := `
		insert into retweet
		(user_id, tweet_id, is_deleted, created_at)
		values
		($1, $2, false, $3)
		on conflict (user_id, tweet_id)
		do update
		set is_deleted = false,
		updated_at = $3
	`

	argsQueryUpsertRetweet := []interface{}{
		param.UserId,
		param.TweetId,
		time.Now(),
	}

	logData["query_upsert_retweet"] = queryUpsertRetweet
	logData["args_query_upsert_retweet"] = fmt.Sprintf("%+v", argsQueryUpsertRetweet)

	errUpsertQuery := repo.db.QueryRowContext(ctx, queryUpsertRetweet, argsQueryUpsertRetweet...).Err()
	if errUpsertQuery != nil {
		logData["error_query_upsert_retweet"] = errUpsertQuery.Error()
		repo.logger.
			WithFields(logData).
			WithError(errUpsertQuery).
			Errorln("error on upsert query")
		span.End()
		return domain.UpsertRetweetResult{}, errUpsertQuery
	}

	span.End()
	return domain.UpsertRetweetResult{}, nil
}
