package app

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	messaging "github.com/DARREN-2000/ai-hypervisor-platform/internal/messaging/nats"
	"github.com/DARREN-2000/ai-hypervisor-platform/internal/storage/postgres"
	rediscache "github.com/DARREN-2000/ai-hypervisor-platform/internal/storage/redis"
)

var (
	connectPostgres = postgres.Connect
	newRedisClient  = rediscache.NewClient
	pingRedis       = rediscache.Ping
	connectNATS     = messaging.Connect
)

// Dependencies holds shared infrastructure clients.
type Dependencies struct {
	Config *config.PlatformConfig
	Logger *logrus.Logger
	DB     *sql.DB
	Redis  *redis.Client
	NATS   *nats.Conn
}

// Build initializes infrastructure dependencies.
func Build(ctx context.Context, cfg *config.PlatformConfig, logger *logrus.Logger) (*Dependencies, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	db, err := connectPostgres(ctx, cfg.Database)
	if err != nil {
		return nil, err
	}

	redisClient := newRedisClient(cfg.Redis)
	if err := pingRedis(ctx, redisClient); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping redis: %w", err)
	}

	natsConn, err := connectNATS(cfg.NATS)
	if err != nil {
		_ = db.Close()
		_ = redisClient.Close()
		return nil, err
	}

	return &Dependencies{
		Config: cfg,
		Logger: logger,
		DB:     db,
		Redis:  redisClient,
		NATS:   natsConn,
	}, nil
}

// Close shuts down infrastructure clients.
func (d *Dependencies) Close(ctx context.Context) error {
	var firstErr error

	if d == nil {
		return nil
	}

	if d.NATS != nil {
		d.NATS.Drain()
		d.NATS.Close()
	}

	if d.Redis != nil {
		if err := d.Redis.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	if d.DB != nil {
		closeCh := make(chan error, 1)
		go func() { closeCh <- d.DB.Close() }()
		select {
		case err := <-closeCh:
			if err != nil && firstErr == nil {
				firstErr = err
			}
		case <-ctx.Done():
			if firstErr == nil {
				firstErr = ctx.Err()
			}
		}
	}

	return firstErr
}
