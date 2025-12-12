package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

func InitRedis() *redis.Client {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}

	return rdb
}

func AddToBlackList(ctx context.Context, rdb *redis.Client, token string, expiration time.Duration) error {
	key := fmt.Sprintf("blacklist:%s", token)

	return rdb.Set(ctx, key, "revoked", expiration).Err()
}

func IsTokenBlackListed(ctx context.Context, rdb *redis.Client, token string) (bool, error) {
	key := fmt.Sprintf("blacklist:%s", token)

	exists, err := rdb.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}
