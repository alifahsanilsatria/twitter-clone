package db

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alifahsanilsatria/twitter-clone/common/constants"
	"github.com/alifahsanilsatria/twitter-clone/domain"
	"github.com/sirupsen/logrus"
)

func (repo *userSessionRepository) GetUserSessionByToken(ctx context.Context, param domain.GetUserSessionByTokenParam) (domain.GetUserSessionByTokenResult, error) {
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
		return domain.GetUserSessionByTokenResult{}, errGetRedis
	}

	return result, nil
}
