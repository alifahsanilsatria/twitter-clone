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

func (repo *tweetRepository) GetTweetByIdAndUserId(ctx context.Context, param domain.GetTweetByIdAndUserIdParam) (domain.GetTweetByIdAndUserIdResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.GetTweetByIdAndUserId", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

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
		span.End()
		return response, errScan
	}

	repo.logger.
		WithFields(logData).
		Infoln("success on GetTweetByIdAndUserId")
	span.End()

	return response, nil
}
