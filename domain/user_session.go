package domain

import (
	"context"
	"time"
)

type SetUserSessionToRedisParam struct {
	UserId int32
	Token  string
	TTL    time.Duration
}

type UserSessionRepository interface {
	SetUserSessionToRedis(ctx context.Context, param SetUserSessionToRedisParam) error
}
