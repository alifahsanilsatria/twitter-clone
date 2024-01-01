package usecase

import (
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

type tweetUsecase struct {
	tweetRepository       domain.TweetRepository
	userSessionRepository domain.UserSessionRepository
	logger                *logrus.Logger
}

func NewTweetUsecase(
	tweetRepository domain.TweetRepository,
	userSessionRepository domain.UserSessionRepository,
	logger *logrus.Logger,
) domain.TweetUsecase {
	return &tweetUsecase{
		tweetRepository:       tweetRepository,
		userSessionRepository: userSessionRepository,
		logger:                logger,
	}
}
