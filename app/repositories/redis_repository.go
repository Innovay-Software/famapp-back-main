package repositories

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	redisClient *redis.Client
	redisCtx    context.Context
}
