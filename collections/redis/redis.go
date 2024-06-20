package redis

import "github.com/redis/go-redis/v9"

var RedisClient *redis.Client

func InitializeRedisClient() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}
