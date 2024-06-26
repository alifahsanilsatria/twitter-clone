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

func (repo *tweetRepository) DeleteLikes(ctx context.Context, param domain.DeleteLikesParam) (domain.DeleteLikesResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.DeleteLikes", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method":     "tweetRepository.DeleteLikes",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	queryDeleteLikes := `
		update likes
		set is_deleted = true,
		updated_at = $3
		where tweet_id = $1
		and user_id = $2
	`

	argsQueryDeleteLikes := []interface{}{
		param.UserId,
		param.TweetId,
		time.Now(),
	}

	logData["query_delete_likes"] = queryDeleteLikes
	logData["args_query_delete_likes"] = fmt.Sprintf("%+v", argsQueryDeleteLikes)

	errDeleteLikesQuery := repo.db.QueryRowContext(ctx, queryDeleteLikes, argsQueryDeleteLikes...).Err()
	if errDeleteLikesQuery != nil {
		logData["error_query_delete_likes"] = errDeleteLikesQuery.Error()
		repo.logger.
			WithFields(logData).
			WithError(errDeleteLikesQuery).
			Errorln("error on delete query")
		span.End()
		return domain.DeleteLikesResult{}, errDeleteLikesQuery
	}

	span.End()
	return domain.DeleteLikesResult{}, nil
}
