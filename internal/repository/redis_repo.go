package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) *RedisRepo {
	return &RedisRepo{client: client}
}

func (r *RedisRepo) Set(ctx context.Context, key string, value string, exp time.Duration) error {
	return r.client.Set(ctx, key, value, exp).Err()
}

func (r *RedisRepo) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisRepo) Del(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}