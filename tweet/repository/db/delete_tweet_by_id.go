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
		"method": "tweetRepository.DeleteTweetById",
		"param":  fmt.Sprintf("%+v", param),
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

	errQueryDeleteTweetById := repo.dbTx.QueryRowContext(ctx, queryDeleteTweetById, argsQueryDeleteTweet...).Err()
	if errQueryDeleteTweetById != nil {
		return errQueryDeleteTweetById
	}

	return nil
}

func (repo *tweetRepository) deleteFromTweetMapChildByTweetId(ctx context.Context, param domain.DeleteTweetByIdParam) error {
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

	errQueryDeleteTweetMapChildByTweetId := repo.dbTx.QueryRowContext(ctx, queryDeleteTweetMapChildByTweetId, argsQueryDeleteTweetMapChildByTweetId...).Err()
	if errQueryDeleteTweetMapChildByTweetId != nil {
		return errQueryDeleteTweetMapChildByTweetId
	}

	return nil
}

func (repo *tweetRepository) deleteFromTweetMapChildByChildTweetId(ctx context.Context, param domain.DeleteTweetByIdParam) error {
	queryDeleteTweetMapChildByTweetId := `
		update tweet_map_child_tweet
		set is_deleted = true,
		updated_at = $1
		where child_tweet_id = $2
	`

	argsQueryDeleteTweetMapChildByTweetId := []interface{}{
		time.Now(),
		param.TweetId,
	}

	errQueryDeleteTweetMapChildByTweetId := repo.dbTx.QueryRowContext(ctx, queryDeleteTweetMapChildByTweetId, argsQueryDeleteTweetMapChildByTweetId...).Err()
	if errQueryDeleteTweetMapChildByTweetId != nil {
		return errQueryDeleteTweetMapChildByTweetId
	}

	return nil
}
