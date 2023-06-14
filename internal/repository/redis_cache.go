package repository

import (
	"github.com/redis/go-redis/v9"
	"os"
	"context"
	"time"

)

type RedisCacheRepository struct {
	Redis  *redis.Client
}

func (r *RedisCacheRepository) InitRedisRepo() *RedisCacheRepository {

	redisHost := os.Getenv("REDIS_HOST")
	rdb := redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	r.Redis = rdb
	return r
}

func(r RedisCacheRepository) SetCacheKey(ctx context.Context, key string, value string, expiration time.Duration) (bool, error) {

	err := r.Redis.Set(ctx, key, value, expiration).Err()

	if err != nil {
		return false ,err
	}

	return true, nil
}


func(r RedisCacheRepository) GetCacheKey(ctx context.Context, key string) (string, error) {

	result, err := r.Redis.Get(ctx, key).Result()

	if err != nil {
		return "" ,err
	}

	return result, nil
}

func(r RedisCacheRepository) DelCacheKey(ctx context.Context, key string) (int64, error) {

	result, err := r.Redis.Del(ctx, key).Result()

	if err != nil {
		return result ,err
	}

	return result, nil
}