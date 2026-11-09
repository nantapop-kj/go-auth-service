package repository

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

var ErrNotFound = errors.New("not found")

type RedisRepository interface {
	Set(ctx context.Context, key string, value string, ttl time.Duration) error
	GetInt(ctx context.Context, key string) (int, error)
	GetDel(ctx context.Context, key string) (string, error)
	Delete(ctx context.Context, key string) error
	IncrWithTTL(ctx context.Context, key string, ttl time.Duration) error
}

type RedisRepo struct {
	client *redis.Client
}

func NewRedisRepo(client *redis.Client) RedisRepository {
	return &RedisRepo{client: client}
}

func (r *RedisRepo) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	return r.client.Set(ctx, key, value, ttl).Err()
}

func (r *RedisRepo) GetInt(ctx context.Context, key string) (int, error) {
	n, err := r.client.Get(ctx, key).Int()
	if errors.Is(err, redis.Nil) {
		return 0, ErrNotFound
	}
	return n, err
}

func (r *RedisRepo) GetDel(ctx context.Context, key string) (string, error) {
	val, err := r.client.GetDel(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", ErrNotFound
	}
	return val, err
}

func (r *RedisRepo) Delete(ctx context.Context, key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisRepo) IncrWithTTL(ctx context.Context, key string, ttl time.Duration) error {
	pipe := r.client.TxPipeline()
	pipe.Incr(ctx, key)
	pipe.ExpireNX(ctx, key, ttl)
	_, err := pipe.Exec(ctx)
	return err
}
