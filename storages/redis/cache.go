package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

func cacheKey(path string) string {
	return fmt.Sprintf("cache:%s", path)
}

func SetCache(ctx context.Context, redisClient *redis.Client, path string, val interface{}) error {
	err := redisClient.Set(ctx, cacheKey(path), val, 5*time.Minute).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetCache(ctx context.Context, redisClient *redis.Client, path string) (string, error) {
	val, err := redisClient.Get(ctx, cacheKey(path)).Result()
	if err != nil {
		return "", err
	}

	return val, nil
}

func DelCache(ctx context.Context, redisClient *redis.Client, path string) error {
	err := redisClient.Del(ctx, cacheKey(path)).Err()
	if err != nil {
		return err
	}

	return nil
}
