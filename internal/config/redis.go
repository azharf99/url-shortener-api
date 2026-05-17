package config

import (
	"context"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

func ConnectRedis() *redis.Client {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}

	password := os.Getenv("REDIS_PASSWORD")
	db := 0 // default db

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := client.Ping(ctx).Result()
	if err != nil {
		zap.L().Warn("Could not connect to Redis, rate limiter will fallback", 
			zap.String("addr", addr), 
			zap.Error(err),
		)
		return nil
	}

	zap.L().Info("Connected to Redis", zap.String("addr", addr))
	return client
}
