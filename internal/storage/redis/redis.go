package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
)

// NewClient creates a Redis client using the supplied configuration.
func NewClient(cfg config.RedisConfig) *redis.Client {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	return redis.NewClient(&redis.Options{
		Addr:         addr,
		Password:     cfg.Password,
		DB:           cfg.Database,
		PoolSize:     cfg.PoolSize,
		DialTimeout:  cfg.ConnectTimeout,
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
	})
}

// Ping validates Redis connectivity.
func Ping(ctx context.Context, client *redis.Client) error {
	return client.Ping(ctx).Err()
}
