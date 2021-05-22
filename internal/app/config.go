package app

import (
	"os"

	redis "github.com/go-redis/redis/v7"
)

//RedisInit initiates redis
func RedisInit() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_DSN"),
	})
	return client
}
