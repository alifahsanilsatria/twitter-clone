package db

import (
	commonWrapper "github.com/alifahsanilsatria/twitter-clone/common/wrapper"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	trace "go.opentelemetry.io/otel/trace"
)

type userRepository struct {
	db     commonWrapper.SQLWrapper
	logger *logrus.Logger
	tracer trace.Tracer
}

func NewUserRepository(
	db commonWrapper.SQLWrapper,
	logger *logrus.Logger,
	tracer trace.Tracer,
) domain.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
		tracer: tracer,
	}
}
