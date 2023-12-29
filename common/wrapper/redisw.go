package wrapper

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisWrapper interface {
	SetEX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Get(ctx context.Context, key string) *redis.StringCmd
}
