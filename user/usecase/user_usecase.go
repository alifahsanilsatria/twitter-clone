package usecase

import (
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

type userUsecase struct {
	userRepository        domain.UserRepository
	userSessionRepository domain.UserSessionRepository
	logger                *logrus.Logger
}

func NewUserUsecase(
	userRepository domain.UserRepository,
	userSessionRepository domain.UserSessionRepository,
	logger *logrus.Logger,
) domain.UserUsecase {
	return &userUsecase{
		userRepository:        userRepository,
		userSessionRepository: userSessionRepository,
		logger:                logger,
	}
}
