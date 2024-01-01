package db

import (
	commonWrapper "github.com/alifahsanilsatria/twitter-clone/common/wrapper"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

type tweetRepository struct {
	db     commonWrapper.SQLWrapper
	dbTx   commonWrapper.SQLTxWrapper
	logger *logrus.Logger
}

func NewUserRepository(
	db commonWrapper.SQLWrapper,
	dbTx commonWrapper.SQLTxWrapper,
	logger *logrus.Logger,
) domain.TweetRepository {
	return &tweetRepository{
		db:     db,
		dbTx:   dbTx,
		logger: logger,
	}
}
