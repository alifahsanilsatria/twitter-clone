package db

import (
	commonWrapper "github.com/alifahsanilsatria/twitter-clone/common/wrapper"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

type userRepository struct {
	db     commonWrapper.SQLWrapper
	logger *logrus.Logger
}

func NewUserRepository(
	db commonWrapper.SQLWrapper,
	logger *logrus.Logger,
) domain.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}
