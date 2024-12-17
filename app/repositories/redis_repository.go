package repositories

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type redisRepo struct {
	redisClient *redis.Client
	redisCtx    context.Context
}

// Set latest file upload time
func (rr *redisRepo) Get(
	key string,
) (
	string, error,
) {
	return rr.redisClient.Get(rr.redisCtx, key).Result()
}

func (rr *redisRepo) Set(
	key, val string,
) error {
	return rr.redisClient.Set(rr.redisCtx, key, val, 0).Err()
}
