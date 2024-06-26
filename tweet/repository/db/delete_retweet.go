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

func (repo *tweetRepository) DeleteRetweet(ctx context.Context, param domain.DeleteRetweetParam) (domain.DeleteRetweetResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.DeleteRetweet", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetRepository.DeleteRetweet",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	queryDeleteRetweet := `
		update retweet
		set is_deleted = true,
		updated_at = $3
		where tweet_id = $1
		and user_id = $2
	`

	argsQueryDeleteRetweet := []interface{}{
		param.TweetId,
		param.UserId,
		time.Now(),
	}

	logData["query_delete_retweet"] = queryDeleteRetweet
	logData["args_query_delete_retweet"] = fmt.Sprintf("%+v", argsQueryDeleteRetweet)

	errDeleteRetweetQuery := repo.dbTx.QueryRowContext(ctx, queryDeleteRetweet, argsQueryDeleteRetweet...).Err()
	if errDeleteRetweetQuery != nil {
		logData["error_query_delete_retweet"] = errDeleteRetweetQuery.Error()
		repo.logger.
			WithFields(logData).
			WithError(errDeleteRetweetQuery).
			Errorln("error on update query")
		span.End()
		return domain.DeleteRetweetResult{}, errDeleteRetweetQuery
	}

	result := domain.DeleteRetweetResult{}
	span.End()
	return result, nil
}
