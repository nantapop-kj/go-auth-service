package repository

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisRepository interface {
	Set(key string, value interface{}, expiration time.Duration) error
	Get(key string) (string, error)
	Delete(key string) error
}

type RedisRepo struct {
	Client *redis.Client
	ctx    context.Context
}

func NewRedisRepo(client *redis.Client) RedisRepository {
	return &RedisRepo{
		Client: client,
		ctx:    context.Background(),
	}
}

func (r *RedisRepo) Set(key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(r.ctx, key, value, expiration).Err()
}

func (r *RedisRepo) Get(key string) (string, error) {
	return r.Client.Get(r.ctx, key).Result()
}

func (r *RedisRepo) Delete(key string) error {
	return r.Client.Del(r.ctx, key).Err()
}
