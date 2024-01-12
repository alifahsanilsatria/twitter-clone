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
		"method":     "tweetRepository.CreateNewTweet",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	var (
		result   domain.CreateNewTweetResult
		errQuery error
	)

	if param.ParentId <= 0 {
		logData["usecase"] = "createNewTweetAsNewTweet"
		result, errQuery = repo.createNewTweetAsNewTweet(ctx, param)
		if errQuery != nil {
			logData["error_query"] = errQuery.Error()
			repo.logger.
				WithFields(logData).
				WithError(errQuery).
				Errorln("error on insert query")
			return domain.CreateNewTweetResult{}, errQuery
		}
	} else {
		logData["usecase"] = "createNewTweetAsReply"
		result, errQuery = repo.createNewTweetAsReply(ctx, param)
		if errQuery != nil {
			logData["error_query"] = errQuery.Error()
			repo.logger.
				WithFields(logData).
				WithError(errQuery).
				Errorln("error on insert query")
			return domain.CreateNewTweetResult{}, errQuery
		}
	}

	repo.logger.
		WithFields(logData).
		Infoln("success CreateNewTweet")

	return result, nil
}

func (repo *tweetRepository) createNewTweetAsNewTweet(ctx context.Context, param domain.CreateNewTweetParam) (domain.CreateNewTweetResult, error) {
	logData := logrus.Fields{
		"method":     "tweetRepository.createNewTweetAsNewTweet",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	queryInsertTweet := `
		insert into tweet
		(user_id, parent_id, content, is_deleted, created_at)
		values
		($1, null, $2, false, $3)
		returning id
	`
	argsQueryInsertTweet := []interface{}{
		param.UserId,
		param.Content,
		time.Now(),
	}

	logData["query_insert_tweet"] = queryInsertTweet
	logData["args_query_insert_tweet"] = fmt.Sprintf("%+v", argsQueryInsertTweet)

	var tweetId int32
	errQuery := repo.db.QueryRowContext(ctx, queryInsertTweet, argsQueryInsertTweet...).Scan(&tweetId)
	if errQuery != nil {
		logData["error_query_insert_tweet"] = errQuery.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQuery).
			Errorln("error on insert query")
		return domain.CreateNewTweetResult{}, errQuery
	}

	repo.logger.
		WithFields(logData).
		Infoln("success createNewTweetAsNewTweet")

	result := domain.CreateNewTweetResult{
		TweetId: tweetId,
	}
	return result, nil
}

func (repo *tweetRepository) createNewTweetAsReply(ctx context.Context, param domain.CreateNewTweetParam) (domain.CreateNewTweetResult, error) {
	logData := logrus.Fields{
		"method":     "tweetRepository.createNewTweetAsReply",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

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

	logData["query_insert_tweet"] = queryInsertTweet
	logData["args_query_insert_tweet"] = fmt.Sprintf("%+v", argsQueryInsertTweet)

	var childTweetId int32
	errQueryInsertTweet := repo.dbTx.QueryRowContext(ctx, queryInsertTweet, argsQueryInsertTweet...).Scan(&childTweetId)
	if errQueryInsertTweet != nil {
		logData["error_query_insert_tweet"] = errQueryInsertTweet.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryInsertTweet).
			Errorln("error on insert query")
		return domain.CreateNewTweetResult{}, errQueryInsertTweet
	}

	logData["child_tweet_id"] = childTweetId

	if childTweetId == 0 {
		errNotGeneratingIdTweet := errors.New("not generating id tweet")
		repo.logger.
			WithFields(logData).
			WithError(errNotGeneratingIdTweet).
			Errorln("error not generating id tweet")
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

	logData["query_insert_tweet_map_child_tweet"] = queryInsertTweetMapChildTweet
	logData["args_query_insert_tweet_map_child_tweet"] = fmt.Sprintf("%+v", argsQueryInsertTweetMapChildTweet)

	errQueryInsertTweetMapChildTweet := repo.dbTx.QueryRowContext(ctx, queryInsertTweetMapChildTweet, argsQueryInsertTweetMapChildTweet...).Err()
	if errQueryInsertTweetMapChildTweet != nil {
		logData["error_query_insert_tweet_map_child_tweet"] = errQueryInsertTweetMapChildTweet.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryInsertTweetMapChildTweet).
			Errorln("error on insert query")
		return domain.CreateNewTweetResult{}, errQueryInsertTweetMapChildTweet
	}

	// Commit the transaction.
	if errCommit := repo.dbTx.Commit(); errCommit != nil {
		logData["error_commit"] = errCommit.Error()
		repo.logger.
			WithFields(logData).
			WithError(errCommit).
			Errorln("error on commit")
		return domain.CreateNewTweetResult{}, errCommit
	}

	repo.logger.
		WithFields(logData).
		Infoln("success createNewTweetAsReply")

	result := domain.CreateNewTweetResult{
		TweetId: childTweetId,
	}

	return result, nil
}
