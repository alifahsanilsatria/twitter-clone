package db

import (
	"context"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *tweetRepository) GetChildrenDataByTweetId(ctx context.Context, param domain.GetChildrenDataByTweetIdParam) (domain.GetChildrenDataByTweetIdResult, error) {
	logData := logrus.Fields{
		"method": "tweetRepository.GetChildrenDataByTweetId",
		"param":  fmt.Sprintf("%+v", param),
	}

	result := domain.GetChildrenDataByTweetIdResult{
		ChildTweet: []domain.GetChildrenDataByTweetIdResult_ChildTweet{},
	}

	query := `
		select t.username, t.complete_name, t.content, 
		coalesce(r.count_retweet, 0) as count_retweet,
		coalesce(l.count_likes, 0) as count_likes,
		coalesce(tmct_child.count_replies, 0) as count_replies
		from tweet_map_child_tweet tmct_parent
		join tweet t
		on tmct_parent.child_tweet_id = t.id
		left join (
			select r.tweet_id, count(r.user_id) as count_retweet
			from retweet r
			where r.is_deleted = false
			group by r.tweet_id
		) as r
		on tmct_parent.child_tweet_id = r.tweet_id
		left join (
			select l.tweet_id, count(l.user_id) as count_likes
			from likes l
			where l.is_deleted = false
			group by l.tweet_id
		) as l
		on tmct_parent.child_tweet_id = l.tweet_id
		left join (
			select tmct_child.tweet_id, count(tmct_child.child_tweet_id) as count_replies
			from tweet_map_child_tweet tmct_child
			where tmct_child.is_deleted = false
			group by tmct_child.tweet_id
		) as tmct_child
		on tmct_parent.child_tweet_id = tmct_child.tweet_id
		where tmct_parent.tweet_id = $1
		and tmct_parent.is_deleted = false
	`

	args := []interface{}{
		param.TweetId,
	}

	queryContextResp, errQuery := repo.db.QueryContext(ctx, query, args...)
	if errQuery != nil {
		logData["error_query"] = errQuery.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQuery).
			Errorln("error on query")
	}

	for queryContextResp.Next() {
		childTweet := domain.GetChildrenDataByTweetIdResult_ChildTweet{}
		queryContextResp.Scan(
			&childTweet.Username,
			&childTweet.CompleteName,
			&childTweet.Content,
			&childTweet.CountRetweet,
			&childTweet.CountLikes,
			&childTweet.CountReplies,
		)
		result.ChildTweet = append(result.ChildTweet, childTweet)
	}

	repo.logger.
		WithFields(logData).
		Infoln("success on GetChildrenDataByTweetId")

	return result, nil
}
