package db

import (
	"context"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/common/constants"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (repo *userSessionRepository) DeleteUserSessionByToken(ctx context.Context, param domain.DeleteUserSessionByTokenParam) error {
	ctx, span := repo.tracer.Start(ctx, "repository.DeleteUserSessionByToken", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method": "userSessionRepository.DeleteUserSessionByToken",
		"param":  fmt.Sprintf("%+v", param),
	}

	redisKeyUserSession := fmt.Sprintf(constants.RedisTemplateKeyUserSession, param.Token)
	redisResult := repo.redisClient.Del(ctx, redisKeyUserSession)

	_, errDelRedis := redisResult.Result()
	if errDelRedis != nil {
		logData["error_del"] = errDelRedis.Error()
		repo.logger.
			WithFields(logData).
			WithError(errDelRedis).
			Errorln("error on del")
		span.End()
		return errDelRedis
	}

	span.End()

	return nil
}
