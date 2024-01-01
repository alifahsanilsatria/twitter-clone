package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/common/constants"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *userSessionRepository) SetUserSessionToRedis(ctx context.Context, param domain.SetUserSessionToRedisParam) error {
	logData := logrus.Fields{
		"method": "userSessionRepository.SetUserSessionToRedis",
		"param":  fmt.Sprintf("%+v", param),
	}

	redisKeyUserSession := fmt.Sprintf(constants.RedisTemplateKeyUserSession, param.Token)
	redisValueUserSession := map[string]interface{}{
		"user_id": param.UserId,
	}
	jsonRedisValueUserSession, _ := json.Marshal(redisValueUserSession)
	result := repo.redisClient.SetEX(ctx, redisKeyUserSession, jsonRedisValueUserSession, param.TTL)
	if result.Err() != nil {
		logData["error_setex"] = result.Err().Error()
		repo.logger.
			WithFields(logData).
			WithError(result.Err()).
			Errorln("error on setex")
		return result.Err()
	}

	return nil
}
