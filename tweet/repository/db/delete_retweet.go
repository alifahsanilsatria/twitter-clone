package db

import (
	"context"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *tweetRepository) DeleteRetweet(ctx context.Context, param domain.DeleteRetweetParam) (domain.DeleteRetweetResult, error) {
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

	errQueryDeleteRetweet := repo.dbTx.QueryRowContext(ctx, queryDeleteRetweet, argsQueryDeleteRetweet...).Err()
	if errQueryDeleteRetweet != nil {
		logData["error_query_delete_retweet"] = errQueryDeleteRetweet.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryDeleteRetweet).
			Errorln("error on update query")
		return domain.DeleteRetweetResult{}, errQueryDeleteRetweet
	}

	result := domain.DeleteRetweetResult{}

	return result, nil
}
