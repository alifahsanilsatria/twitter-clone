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

type GetUserSessionByTokenParam struct {
	Token string
}

type GetUserSessionByTokenResult struct {
	UserId int32 `json:"user_id"`
}

type DeleteUserSessionByTokenParam struct {
	Token string
}

type UserSessionRepository interface {
	SetUserSessionToRedis(ctx context.Context, param SetUserSessionToRedisParam) error
	GetUserSessionByToken(ctx context.Context, param GetUserSessionByTokenParam) (GetUserSessionByTokenResult, error)
	DeleteUserSessionByToken(ctx context.Context, param DeleteUserSessionByTokenParam) error
}
