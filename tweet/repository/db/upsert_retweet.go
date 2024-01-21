package db

import (
	"context"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *tweetRepository) UpsertRetweet(ctx context.Context, param domain.UpsertRetweetParam) (domain.UpsertRetweetResult, error) {
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
		return domain.UpsertRetweetResult{}, errUpsertQuery
	}

	return domain.UpsertRetweetResult{}, nil
}
