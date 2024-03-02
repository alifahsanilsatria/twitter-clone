package usecase

import (
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	trace "go.opentelemetry.io/otel/trace"
)

type userUsecase struct {
	userRepository        domain.UserRepository
	userSessionRepository domain.UserSessionRepository
	logger                *logrus.Logger
	tracer                trace.Tracer
}

func NewUserUsecase(
	userRepository domain.UserRepository,
	userSessionRepository domain.UserSessionRepository,
	logger *logrus.Logger,
	tracer trace.Tracer,
) domain.UserUsecase {
	return &userUsecase{
		userRepository:        userRepository,
		userSessionRepository: userSessionRepository,
		logger:                logger,
		tracer:                tracer,
	}
}
