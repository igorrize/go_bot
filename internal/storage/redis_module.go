package storage

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
)

var ctx = context.Background()

var RedisClient *redis.Client

func InitRedisClient() *redis.Client {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "",
		DB:       0,
	})
	return RedisClient
}

func SetKey(key, value string) error {
	if RedisClient == nil {
		return redis.Nil
	}

	err := RedisClient.Set(ctx, key, value, 0).Err()
	return err
}

func GetKey(key string) string {
	if RedisClient == nil {
		return "no client"
	}

	val, _ := RedisClient.Get(ctx, key).Result()
	return val
}

func CloseRedisClient(*redis.Client) {
	if RedisClient != nil {
		err := RedisClient.Close()
		if err != nil {
			log.Println("Error closing Redis client:", err)
		}
	}
}
