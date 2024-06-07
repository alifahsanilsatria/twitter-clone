package usecase

import (
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	trace "go.opentelemetry.io/otel/trace"
)

type tweetUsecase struct {
	tweetRepository       domain.TweetRepository
	userRepository        domain.UserRepository
	userSessionRepository domain.UserSessionRepository
	logger                *logrus.Logger
	tracer                trace.Tracer
}

func NewTweetUsecase(
	tweetRepository domain.TweetRepository,
	userRepository domain.UserRepository,
	userSessionRepository domain.UserSessionRepository,
	logger *logrus.Logger,
	tracer trace.Tracer,
) domain.TweetUsecase {
	return &tweetUsecase{
		tweetRepository:       tweetRepository,
		userRepository:        userRepository,
		userSessionRepository: userSessionRepository,
		logger:                logger,
		tracer:                tracer,
	}
}
