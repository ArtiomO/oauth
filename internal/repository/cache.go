package repository

import (
	"context"
	"time"
)

type CacheRepository interface {
	SetCacheKey(ctx context.Context, key string, value string, expiration time.Duration) (bool, error)
	GetCacheKey(ctx context.Context, key string) (string, error)
	DelCacheKey(ctx context.Context, key string) (int64, error)
}