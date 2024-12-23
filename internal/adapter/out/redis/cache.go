package redis

import (
	"context"
	"time"

	"github.com/VanillaSkys/todo_fiber/internal/core/port/cache"
	"github.com/redis/go-redis/v9"
)

type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type redisCache struct {
	client *redis.Client
}

func NewRedisCache(client *redis.Client) cache.Cache {
	return &redisCache{client: client}
}

func (r *redisCache) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *redisCache) Set(ctx context.Context, key string, value string, expiration int64) error {
	duration := time.Duration(expiration)
	return r.client.Set(ctx, key, value, duration).Err()
}

func (r *redisCache) Del(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	return err
}
