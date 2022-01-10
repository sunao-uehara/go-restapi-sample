package redis

import (
	"github.com/go-redis/redis/v8"
)

func Initialize(url string) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	return rdb
}
