package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *tweetRepository) CreateNewTweet(ctx context.Context, param domain.CreateNewTweetParam) (domain.CreateNewTweetResult, error) {
	logData := logrus.Fields{
		"method": "tweetRepository.CreateNewTweet",
		"param":  fmt.Sprintf("%+v", param),
	}

	var (
		result   domain.CreateNewTweetResult
		errQuery error
	)

	if param.ParentId <= 0 {
		result, errQuery = repo.createNewTweetAsNewTweet(ctx, param)
		if errQuery != nil {
			logData["error_query"] = errQuery.Error()
			logData["usecase"] = "createNewTweetAsNewTweet"
			repo.logger.
				WithFields(logData).
				WithError(errQuery).
				Errorln("error on insert query")
			return domain.CreateNewTweetResult{}, errQuery
		}
	} else {
		result, errQuery = repo.createNewTweetAsReply(ctx, param)
		if errQuery != nil {
			logData["error_query"] = errQuery.Error()
			logData["usecase"] = "createNewTweetAsReply"
			repo.logger.
				WithFields(logData).
				WithError(errQuery).
				Errorln("error on insert query")
			return domain.CreateNewTweetResult{}, errQuery
		}
	}

	return result, nil
}

func (repo *tweetRepository) createNewTweetAsNewTweet(ctx context.Context, param domain.CreateNewTweetParam) (domain.CreateNewTweetResult, error) {
	query := `
		insert into tweet
		(user_id, parent_id, content, is_deleted, created_at)
		values
		($1, null, $2, false, $3)
		returning id
	`
	args := []interface{}{
		param.UserId,
		param.Content,
		time.Now(),
	}

	var tweetId int32
	errQuery := repo.db.QueryRowContext(ctx, query, args...).Scan(&tweetId)
	if errQuery != nil {
		return domain.CreateNewTweetResult{}, errQuery
	}

	result := domain.CreateNewTweetResult{
		TweetId: tweetId,
	}
	return result, nil
}

func (repo *tweetRepository) createNewTweetAsReply(ctx context.Context, param domain.CreateNewTweetParam) (domain.CreateNewTweetResult, error) {
	// Defer a rollback in case anything fails.
	defer repo.dbTx.Rollback()

	queryInsertTweet := `
		insert into tweet
		(user_id, parent_id, content, is_deleted, created_at)
		values
		($1, $2, $3, false, $4)
		returning id
	`
	argsQueryInsertTweet := []interface{}{
		param.UserId,
		param.ParentId,
		param.Content,
		time.Now(),
	}

	var childTweetId int32
	errQueryInsertTweet := repo.dbTx.QueryRowContext(ctx, queryInsertTweet, argsQueryInsertTweet...).Scan(&childTweetId)
	if errQueryInsertTweet != nil {
		return domain.CreateNewTweetResult{}, errQueryInsertTweet
	}

	if childTweetId == 0 {
		return domain.CreateNewTweetResult{}, errors.New("not generating id tweet")
	}

	queryInsertTweetMapChildTweet := `
		insert into tweet_map_child_tweet
		(tweet_id, child_tweet_id, is_deleted, created_at)
		values
		($1, $2, false, $3)
	`

	argsQueryInsertTweetMapChildTweet := []interface{}{
		param.ParentId,
		childTweetId,
		time.Now(),
	}

	errQueryInsertTweetMapChildTweet := repo.dbTx.QueryRowContext(ctx, queryInsertTweetMapChildTweet, argsQueryInsertTweetMapChildTweet...).Err()
	if errQueryInsertTweetMapChildTweet != nil {
		return domain.CreateNewTweetResult{}, errQueryInsertTweetMapChildTweet
	}

	// Commit the transaction.
	if errCommit := repo.dbTx.Commit(); errCommit != nil {
		return domain.CreateNewTweetResult{}, errCommit
	}

	result := domain.CreateNewTweetResult{
		TweetId: childTweetId,
	}

	return result, nil
}
