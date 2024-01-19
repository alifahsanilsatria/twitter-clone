package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *tweetRepository) GetTweetByIdAndUserId(ctx context.Context, param domain.GetTweetByIdAndUserIdParam) (domain.GetTweetByIdAndUserIdResult, error) {
	logData := logrus.Fields{
		"method": "tweetRepository.GetTweetByIdAndUserId",
		"param":  fmt.Sprintf("%+v", param),
	}

	query := `
		select id 
		from tweet
		where id = $1
		and user_id = $2
		and is_deleted = false
	`

	args := []interface{}{
		param.TweetId,
		param.UserId,
	}

	queryRowContextResp := repo.db.QueryRowContext(ctx, query, args...)

	response := domain.GetTweetByIdAndUserIdResult{}
	errScan := queryRowContextResp.Scan(&response.TweetId)
	if errScan != nil && errScan != sql.ErrNoRows {
		logData["error_scan"] = errScan.Error()
		repo.logger.
			WithFields(logData).
			WithError(errScan).
			Errorln("error on scan")
		return response, errScan
	}

	repo.logger.
		WithFields(logData).
		Infoln("success on GetTweetByIdAndUserId")

	return response, nil
}
