package db

import (
	commonWrapper "github.com/alifahsanilsatria/twitter-clone/common/wrapper"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	trace "go.opentelemetry.io/otel/trace"
)

type tweetRepository struct {
	db     commonWrapper.SQLWrapper
	dbTx   commonWrapper.SQLTxWrapper
	logger *logrus.Logger
	tracer trace.Tracer
}

func NewTweetRepository(
	db commonWrapper.SQLWrapper,
	dbTx commonWrapper.SQLTxWrapper,
	logger *logrus.Logger,
	tracer trace.Tracer,
) domain.TweetRepository {
	return &tweetRepository{
		db:     db,
		dbTx:   dbTx,
		logger: logger,
		tracer: tracer,
	}
}
