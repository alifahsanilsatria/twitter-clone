package db

import (
	"context"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (repo *tweetRepository) GetParentsDataByTweetId(ctx context.Context, param domain.GetParentsDataByTweetIdParam) (domain.GetParentsDataByTweetIdResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.GetParentsDataByTweetId", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method": "tweetRepository.GetParentsDataByTweetId",
		"param":  fmt.Sprintf("%+v", param),
	}

	result := domain.GetParentsDataByTweetIdResult{
		Parent: []domain.GetParentsDataByTweetIdResult_Parent{},
	}

	query := `
		with recursive anchestor as (
			select t.id, u.username, u.complete_name, t.content, 
			coalesce(r.count_retweet, 0) as count_retweet,
			coalesce(l.count_likes, 0) as count_likes,
			coalesce(tmct_child.count_replies, 0) as count_replies
			from tweet_map_child_tweet tmct_parent
			join tweet t
			on tmct_parent.tweet_id = t.id
			join users u
			on t.user_id = u.id
			left join (
				select r.tweet_id, count(r.user_id) as count_retweet
				from retweet r
				where r.is_deleted = false
				group by r.tweet_id
			) as r
			on tmct_parent.tweet_id = r.tweet_id
			left join (
				select l.tweet_id, count(l.user_id) as count_likes
				from likes l
				where l.is_deleted = false
				group by l.tweet_id
			) as l
			on tmct_parent.tweet_id = l.tweet_id
			left join (
				select tmct_child.tweet_id, count(tmct_child.child_tweet_id) as count_replies
				from tweet_map_child_tweet tmct_child
				where tmct_child.is_deleted = false
				group by tmct_child.tweet_id
			) as tmct_child
			on tmct_parent.tweet_id = tmct_child.tweet_id
			where tmct_parent.child_tweet_id = $1
			and tmct_parent.is_deleted = false
			union all
			select t.id, u.username, u.complete_name, t.content, 
			coalesce(r.count_retweet, 0) as count_retweet,
			coalesce(l.count_likes, 0) as count_likes,
			coalesce(tmct_child.count_replies, 0) as count_replies
			from tweet_map_child_tweet tmct_parent
			join anchestor a
			on tmct_parent.child_tweet_id = a.id
			join tweet t
			on tmct_parent.tweet_id = t.id
			join users u
			on t.user_id = u.id
			left join (
				select r.tweet_id, count(r.user_id) as count_retweet
				from retweet r
				where r.is_deleted = false
				group by r.tweet_id
			) as r
			on tmct_parent.tweet_id = r.tweet_id
			left join (
				select l.tweet_id, count(l.user_id) as count_likes
				from likes l
				where l.is_deleted = false
				group by l.tweet_id
			) as l
			on tmct_parent.tweet_id = l.tweet_id
			left join (
				select tmct_child.tweet_id, count(tmct_child.child_tweet_id) as count_replies
				from tweet_map_child_tweet tmct_child
				where tmct_child.is_deleted = false
				group by tmct_child.tweet_id
			) as tmct_child
			on tmct_parent.tweet_id = tmct_child.tweet_id
			where tmct_parent.is_deleted = false
		)
		select 
			id, username, complete_name, content, 
			count_retweet, count_likes, count_replies
		from anchestor;
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
		span.End()
	}

	for queryContextResp.Next() {
		parent := domain.GetParentsDataByTweetIdResult_Parent{}
		queryContextResp.Scan(
			&parent.TweetId,
			&parent.Username,
			&parent.CompleteName,
			&parent.Content,
			&parent.CountRetweet,
			&parent.CountLikes,
			&parent.CountReplies,
		)
		result.Parent = append(result.Parent, parent)
	}

	finalSlicesOfParent := make([]domain.GetParentsDataByTweetIdResult_Parent, len(result.Parent))
	for idx, _ := range result.Parent {
		finalSlicesOfParent[idx] = result.Parent[len(result.Parent)-1-idx]
	}

	result.Parent = finalSlicesOfParent

	repo.logger.
		WithFields(logData).
		Infoln("success on GetChildrenDataByTweetId")
	span.End()

	return result, nil
}
