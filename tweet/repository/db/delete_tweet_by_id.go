package db

import (
	"context"
	"fmt"
	"time"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *tweetRepository) DeleteTweetById(ctx context.Context, param domain.DeleteTweetByIdParam) (domain.DeleteTweetByIdResult, error) {
	logData := logrus.Fields{
		"method":     "tweetRepository.DeleteTweetById",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	defer repo.dbTx.Rollback()

	errQueryDeleteFromTweetTable := repo.deleteFromTweetTableById(ctx, param)
	if errQueryDeleteFromTweetTable != nil {
		logData["error_delete_from_tweet_table"] = errQueryDeleteFromTweetTable.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryDeleteFromTweetTable).
			Errorln("error on deleteFromTweetTableById")
		return domain.DeleteTweetByIdResult{}, errQueryDeleteFromTweetTable
	}

	errQueryDeleteTweetMapChildByTweetId := repo.deleteFromTweetMapChildByTweetId(ctx, param)
	if errQueryDeleteTweetMapChildByTweetId != nil {
		logData["error_delete_from_tweet_map_child_by_tweet_id"] = errQueryDeleteTweetMapChildByTweetId.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryDeleteTweetMapChildByTweetId).
			Errorln("error on deleteFromTweetMapChildByTweetId")
		return domain.DeleteTweetByIdResult{}, errQueryDeleteTweetMapChildByTweetId
	}

	errQueryDeleteTweetMapChildByChildTweetId := repo.deleteFromTweetMapChildByChildTweetId(ctx, param)
	if errQueryDeleteTweetMapChildByChildTweetId != nil {
		logData["error_delete_from_tweet_map_child_by_child_tweet_id"] = errQueryDeleteTweetMapChildByChildTweetId.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryDeleteTweetMapChildByChildTweetId).
			Errorln("error on deleteFromTweetMapChildByChildTweetId")
		return domain.DeleteTweetByIdResult{}, errQueryDeleteTweetMapChildByChildTweetId
	}

	// Commit the transaction.
	if errCommit := repo.dbTx.Commit(); errCommit != nil {
		logData["error_commit"] = errCommit.Error()
		repo.logger.
			WithFields(logData).
			WithError(errCommit).
			Errorln("error on deleteFromTweetMapChildByChildTweetId")
		return domain.DeleteTweetByIdResult{}, errCommit
	}

	repo.logger.
		WithFields(logData).
		Infoln("success on DeleteTweetById")

	result := domain.DeleteTweetByIdResult{
		TweetId: param.TweetId,
	}

	return result, nil
}

func (repo *tweetRepository) deleteFromTweetTableById(ctx context.Context, param domain.DeleteTweetByIdParam) error {
	logData := logrus.Fields{
		"method":     "tweetRepository.deleteFromTweetTableById",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	queryDeleteTweetById := `
		update tweet
		set is_deleted = true,
		updated_at = $1
		where id = $2
	`

	argsQueryDeleteTweet := []interface{}{
		time.Now(),
		param.TweetId,
	}

	logData["query_delete_tweet_by_id"] = queryDeleteTweetById
	logData["args_query_delete_tweet_by_id"] = fmt.Sprintf("%+v", argsQueryDeleteTweet)

	errQueryDeleteTweetById := repo.dbTx.QueryRowContext(ctx, queryDeleteTweetById, argsQueryDeleteTweet...).Err()
	if errQueryDeleteTweetById != nil {
		logData["error_query_delete_tweet_by_id"] = errQueryDeleteTweetById.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryDeleteTweetById).
			Errorln("error on deleteFromTweetTableById")
		return errQueryDeleteTweetById
	}

	repo.logger.
		WithFields(logData).
		Infoln("success deleteFromTweetTableById")

	return nil
}

func (repo *tweetRepository) deleteFromTweetMapChildByTweetId(ctx context.Context, param domain.DeleteTweetByIdParam) error {
	logData := logrus.Fields{
		"method":     "tweetRepository.deleteFromTweetTableById",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	queryDeleteTweetMapChildByTweetId := `
		update tweet_map_child_tweet
		set is_deleted = true,
		updated_at = $1
		where tweet_id = $2
	`

	argsQueryDeleteTweetMapChildByTweetId := []interface{}{
		time.Now(),
		param.TweetId,
	}

	logData["query_delete_tweet_map_child_by_tweet_id"] = queryDeleteTweetMapChildByTweetId
	logData["args_query_delete_tweet_map_child_by_tweet_id"] = fmt.Sprintf("%+v", argsQueryDeleteTweetMapChildByTweetId)

	errQueryDeleteTweetMapChildByTweetId := repo.dbTx.QueryRowContext(ctx, queryDeleteTweetMapChildByTweetId, argsQueryDeleteTweetMapChildByTweetId...).Err()
	if errQueryDeleteTweetMapChildByTweetId != nil {
		logData["error_query_delete_tweet_map_child_by_tweet_id"] = errQueryDeleteTweetMapChildByTweetId.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryDeleteTweetMapChildByTweetId).
			Errorln("error on deleteFromTweetMapChildByTweetId")
		return errQueryDeleteTweetMapChildByTweetId
	}

	repo.logger.
		WithFields(logData).
		Infoln("success deleteFromTweetMapChildByTweetId")

	return nil
}

func (repo *tweetRepository) deleteFromTweetMapChildByChildTweetId(ctx context.Context, param domain.DeleteTweetByIdParam) error {
	logData := logrus.Fields{
		"method":     "tweetRepository.deleteFromTweetMapChildByChildTweetId",
		"request_id": ctx.Value("request_id"),
		"param":      fmt.Sprintf("%+v", param),
	}

	queryDeleteTweetMapChildByChildTweetId := `
		update tweet_map_child_tweet
		set is_deleted = true,
		updated_at = $1
		where child_tweet_id = $2
	`

	argsQueryDeleteTweetMapChildByChildTweetId := []interface{}{
		time.Now(),
		param.TweetId,
	}

	logData["query_delete_tweet_map_child_by_child_tweet_id"] = queryDeleteTweetMapChildByChildTweetId
	logData["args_query_delete_tweet_map_child_by_child_tweet_id"] = fmt.Sprintf("%+v", argsQueryDeleteTweetMapChildByChildTweetId)

	errQueryDeleteTweetMapChildByChildTweetId := repo.dbTx.QueryRowContext(ctx, queryDeleteTweetMapChildByChildTweetId, argsQueryDeleteTweetMapChildByChildTweetId...).Err()
	if errQueryDeleteTweetMapChildByChildTweetId != nil {
		logData["error_query_delete_tweet_map_child_by_child_tweet_id"] = errQueryDeleteTweetMapChildByChildTweetId.Error()
		repo.logger.
			WithFields(logData).
			WithError(errQueryDeleteTweetMapChildByChildTweetId).
			Errorln("error on deleteFromTweetMapChildByTweetId")
		return errQueryDeleteTweetMapChildByChildTweetId
	}

	repo.logger.
		WithFields(logData).
		Infoln("success deleteFromTweetMapChildByChildTweetId")

	return nil
}
