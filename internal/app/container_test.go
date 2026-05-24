package app

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"io"
	"sync/atomic"
	"testing"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

type stubDriver struct {
	closeCount *atomic.Int32
}

type stubConn struct {
	closeCount *atomic.Int32
}

func (d stubDriver) Open(string) (driver.Conn, error) {
	return &stubConn{closeCount: d.closeCount}, nil
}

func (c *stubConn) Prepare(string) (driver.Stmt, error) { return stubStmt{}, nil }
func (c *stubConn) Close() error {
	c.closeCount.Add(1)
	return nil
}
func (c *stubConn) Begin() (driver.Tx, error) { return stubTx{}, nil }

type stubStmt struct{}

func (stubStmt) Close() error                                        { return nil }
func (stubStmt) NumInput() int                                        { return 0 }
func (stubStmt) Exec([]driver.Value) (driver.Result, error)           { return stubResult{}, nil }
func (stubStmt) Query([]driver.Value) (driver.Rows, error)            { return stubRows{}, nil }

type stubTx struct{}

func (stubTx) Commit() error   { return nil }
func (stubTx) Rollback() error  { return nil }

type stubResult struct{}

func (stubResult) LastInsertId() (int64, error) { return 0, nil }
func (stubResult) RowsAffected() (int64, error) { return 0, nil }

type stubRows struct{}

func (stubRows) Columns() []string               { return nil }
func (stubRows) Close() error                    { return nil }
func (stubRows) Next([]driver.Value) error       { return io.EOF }

func TestBuildReturnsPostgresError(t *testing.T) {
	origConnectPostgres := connectPostgres
	origNewRedisClient := newRedisClient
	origPingRedis := pingRedis
	origConnectNATS := connectNATS
	defer func() {
		connectPostgres = origConnectPostgres
		newRedisClient = origNewRedisClient
		pingRedis = origPingRedis
		connectNATS = origConnectNATS
	}()

	connectPostgres = func(context.Context, config.DatabaseConfig) (*sql.DB, error) {
		return nil, errors.New("postgres unavailable")
	}
	newRedisClient = func(config.RedisConfig) *redis.Client {
		t.Fatal("redis client should not be created when postgres fails")
		return nil
	}
	pingRedis = func(context.Context, *redis.Client) error {
		t.Fatal("redis ping should not run when postgres fails")
		return nil
	}
	connectNATS = func(config.NATSConfig) (*nats.Conn, error) {
		t.Fatal("nats should not be created when postgres fails")
		return nil, nil
	}

	_, err := Build(context.Background(), config.DefaultConfig(), logrus.New())
	if err == nil || err.Error() != "postgres unavailable" {
		t.Fatalf("expected postgres unavailable error, got %v", err)
	}
}

func TestBuildClosesDatabaseWhenRedisPingFails(t *testing.T) {
	origConnectPostgres := connectPostgres
	origNewRedisClient := newRedisClient
	origPingRedis := pingRedis
	origConnectNATS := connectNATS
	defer func() {
		connectPostgres = origConnectPostgres
		newRedisClient = origNewRedisClient
		pingRedis = origPingRedis
		connectNATS = origConnectNATS
	}()

	closeCount := &atomic.Int32{}
	sql.Register("stub-app-close", stubDriver{closeCount: closeCount})

	connectPostgres = func(context.Context, config.DatabaseConfig) (*sql.DB, error) {
		db, err := sql.Open("stub-app-close", "")
		if err != nil {
			return nil, err
		}
		if err := db.Ping(); err != nil {
			return nil, err
		}
		return db, nil
	}
	newRedisClient = func(config.RedisConfig) *redis.Client {
		return redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	}
	pingRedis = func(context.Context, *redis.Client) error {
		return errors.New("redis unavailable")
	}
	connectNATS = func(config.NATSConfig) (*nats.Conn, error) {
		t.Fatal("nats should not be created when redis ping fails")
		return nil, nil
	}

	_, err := Build(context.Background(), config.DefaultConfig(), logrus.New())
	if err == nil || err.Error() != "ping redis: redis unavailable" {
		t.Fatalf("expected redis ping error, got %v", err)
	}
	if got := closeCount.Load(); got != 1 {
		t.Fatalf("expected database close to be called once, got %d", got)
	}
}
