package db

import (
	commonWrapper "github.com/alifahsanilsatria/twitter-clone/common/wrapper"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	trace "go.opentelemetry.io/otel/trace"
)

type userSessionRepository struct {
	redisClient commonWrapper.RedisWrapper
	logger      *logrus.Logger
	tracer      trace.Tracer
}

func NewUserSessionRepository(
	redisClient commonWrapper.RedisWrapper,
	logger *logrus.Logger,
	tracer trace.Tracer,
) domain.UserSessionRepository {
	return &userSessionRepository{
		redisClient: redisClient,
		logger:      logger,
		tracer:      tracer,
	}
}
