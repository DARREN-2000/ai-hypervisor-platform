package nats

import (
	"fmt"
	"strings"
	"time"

	"github.com/nats-io/nats.go"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
)

// Connect establishes a NATS connection with the supplied configuration.
func Connect(cfg config.NATSConfig) (*nats.Conn, error) {
	urls := strings.Join(cfg.URLs, ",")
	options := []nats.Option{
		nats.Timeout(timeoutOrDefault(cfg.ConnectTimeout, 5*time.Second)),
		nats.MaxReconnects(cfg.MaxReconnect),
		nats.ReconnectWait(timeoutOrDefault(cfg.ReconnectWait, 2*time.Second)),
		nats.Name("ai-hypervisor-platform"),
	}

	if cfg.Username != "" || cfg.Password != "" {
		options = append(options, nats.UserInfo(cfg.Username, cfg.Password))
	}

	conn, err := nats.Connect(urls, options...)
	if err != nil {
		return nil, fmt.Errorf("connect nats: %w", err)
	}

	return conn, nil
}

func timeoutOrDefault(value, fallback time.Duration) time.Duration {
	if value > 0 {
		return value
	}
	return fallback
}
