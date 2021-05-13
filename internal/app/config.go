package app

import (
	"os"

	redis "github.com/go-redis/redis/v7"
)

func RedisInit() *redis.Client {
	var client *redis.Client
	client = redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_DSN"),
	})
	return client
}
