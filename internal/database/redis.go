package database

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/kodokbakar/go-ecommerce-api/internal/config"
)

func NewRedisOptions(cfg config.RedisConfig) *redis.Options {
	return &redis.Options{
		Addr:         cfg.Addr(),
		Password:     cfg.Password,
		DB:           cfg.DB,
		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,
		DialTimeout:  cfg.DialTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	}
}

func NewRedisClient(ctx context.Context, cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(NewRedisOptions(cfg))

	if err := client.Ping(ctx).Err(); err != nil {
		if closeErr := client.Close(); closeErr != nil {
			return nil, fmt.Errorf("failed to ping redis: %w; also failed to close redis client: %v", err, closeErr)
		}

		return nil, fmt.Errorf("failed to ping redis: %w", err)
	}

	return client, nil
}
