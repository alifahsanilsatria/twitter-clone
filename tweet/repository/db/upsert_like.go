package db

import (
	"context"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *tweetRepository) UpsertLikes(ctx context.Context, param domain.UpsertLikesParam) (domain.UpsertLikesResult, error) {
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
		do update likes
		set is_deleted = false,
		updated_at = $3
		where user_id = $1
		and tweet_id = $2
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
		return domain.UpsertLikesResult{}, errUpsertQuery
	}

	return domain.UpsertLikesResult{}, nil
}
