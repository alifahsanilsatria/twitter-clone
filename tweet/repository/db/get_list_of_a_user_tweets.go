package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (repo *tweetRepository) GetListOfAUserTweets(ctx context.Context, param domain.GetListOfAUserTweetsParam) (domain.GetListOfAUserTweetsResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.GetListOfAUserTweets", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method": "tweetRepository.GetListOfAUserTweets",
		"param":  fmt.Sprintf("%+v", param),
	}

	query := `
		select t.id, u.username, u.complete_name, t.content,
		coalesce(r.count_retweet, 0) as count_retweet,
		coalesce(l.count_likes, 0) as count_likes,
		coalesce(tmct.count_replies, 0) as count_replies
		from tweet t
		join users u on t.user_id = u.id
		left join (
		select r.tweet_id, count(r.user_id) as count_retweet
		from retweet r
		where r.is_deleted = false
		group by r.tweet_id
		) r on t.id = r.tweet_id
		left join (
		select l.tweet_id, count(l.user_id) as count_likes
		from likes l
		where l.is_deleted = false
		group by l.tweet_id
		) l on t.id = l.tweet_id
		left join (
		select tmct.tweet_id, count(tmct.child_tweet_id) as count_replies
		from tweet_map_child_tweet tmct
		where tmct.is_deleted = false
		group by tmct.tweet_id
		) tmct on t.id = tmct.tweet_id
		where u.id = $1 AND t.is_deleted = false and u.is_deleted = false;
	`

	args := []interface{}{
		param.UserId,
	}

	queryContextResp, errQuery := repo.db.QueryContext(ctx, query, args...)
	if errQuery != nil {
		return domain.GetListOfAUserTweetsResult{}, errQuery
	}

	result := domain.GetListOfAUserTweetsResult{}

	for queryContextResp.NextResultSet() {
		tweet := domain.GetListOfAUserTweetsResult_Tweet{}
		errScan := queryContextResp.Scan(
			&tweet.TweetId,
			&tweet.Username,
			&tweet.CompleteName,
			&tweet.Content,
			&tweet.CountRetweet,
			&tweet.CountLikes,
			&tweet.CountReplies,
		)
		if errScan != nil && errScan != sql.ErrNoRows {
			logData["error_scan"] = errScan.Error()
			repo.logger.
				WithFields(logData).
				WithError(errScan).
				Errorln("error on scan")
			span.End()
			return domain.GetListOfAUserTweetsResult{}, errScan
		}

		result.Tweets = append(result.Tweets, tweet)
	}

	repo.logger.
		WithFields(logData).
		Infoln("success on GetListOfAUserTweets")
	span.End()

	return result, nil
}
