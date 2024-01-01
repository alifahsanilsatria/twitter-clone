package db

import (
	commonWrapper "github.com/alifahsanilsatria/twitter-clone/common/wrapper"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

type userSessionRepository struct {
	redisClient commonWrapper.RedisWrapper
	logger      *logrus.Logger
}

func NewUserSessionRepository(
	redisClient commonWrapper.RedisWrapper,
	logger *logrus.Logger,
) domain.UserSessionRepository {
	return &userSessionRepository{
		redisClient: redisClient,
		logger:      logger,
	}
}
