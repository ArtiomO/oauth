package repository

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"os"
	"time"
)

type RedisCacheRepository struct {
	redis *redis.Client
}

func (r *RedisCacheRepository) InitRedisRepo() *RedisCacheRepository {

	redisHost := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r.redis = rdb
	return r
}

func (r RedisCacheRepository) Disconnect() (bool, error) {

	err := r.redis.Close()

	log.Printf("Disconnecting redis.")

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r RedisCacheRepository) SetCacheKey(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {

	err := r.redis.Set(ctx, key, value, expiration).Err()

	if err != nil {
		return false, err
	}

	return true, nil
}

func (r RedisCacheRepository) GetCacheKey(ctx context.Context, key string) (string, error) {

	result, err := r.redis.Get(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return "", ErrKeyDoesntExists
	}

	return result, nil
}

func (r RedisCacheRepository) DelCacheKey(ctx context.Context, key string) (int64, error) {

	result, err := r.redis.Del(ctx, key).Result()

	if errors.Is(err, redis.Nil) {
		return 0, ErrKeyDoesntExists
	}

	return result, nil
}
