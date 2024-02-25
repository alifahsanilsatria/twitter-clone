package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/common/constants"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (repo *userSessionRepository) GetUserSessionByToken(ctx context.Context, param domain.GetUserSessionByTokenParam) (domain.GetUserSessionByTokenResult, error) {
	ctx, span := repo.tracer.Start(ctx, "repository.GetUserSessionByToken", trace.WithAttributes(
		attribute.String("param", fmt.Sprintf("%+v", param)),
	))

	logData := logrus.Fields{
		"method": "userSessionRepository.GetUserSessionByToken",
		"param":  fmt.Sprintf("%+v", param),
	}

	redisKeyUserSession := fmt.Sprintf(constants.RedisTemplateKeyUserSession, param.Token)
	redisResult := repo.redisClient.Get(ctx, redisKeyUserSession)

	resultStr, errGetRedis := redisResult.Result()
	if errGetRedis != nil {
		logData["error_get"] = errGetRedis.Error()
		repo.logger.
			WithFields(logData).
			WithError(errGetRedis).
			Errorln("error on get")
		span.End()
		return domain.GetUserSessionByTokenResult{}, errGetRedis
	}

	var result domain.GetUserSessionByTokenResult
	errUnmarshal := json.Unmarshal([]byte(resultStr), &result)
	if errUnmarshal != nil {
		logData["error_unmarshal"] = errUnmarshal.Error()
		repo.logger.
			WithFields(logData).
			WithError(errUnmarshal).
			Errorln("error on unmarshal")
		span.End()
		return domain.GetUserSessionByTokenResult{}, errGetRedis
	}
	span.End()

	return result, nil
}
