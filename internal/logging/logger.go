package logging

import (
	"io"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/DARREN-2000/ai-hypervisor-platform/internal/config"
)

// NewLogger creates a structured logger for a service.
func NewLogger(service string, cfg config.LoggingConfig) *logrus.Logger {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{TimestampFormat: time.RFC3339})
	logger.SetLevel(parseLevel(cfg.Level))
	logger.SetOutput(resolveOutput(cfg.OutputPath))

	return logger
}

func parseLevel(level string) logrus.Level {
	switch strings.ToLower(level) {
	case "debug":
		return logrus.DebugLevel
	case "warn", "warning":
		return logrus.WarnLevel
	case "error":
		return logrus.ErrorLevel
	default:
		return logrus.InfoLevel
	}
}

func resolveOutput(path string) io.Writer {
	switch strings.ToLower(path) {
	case "", "stdout":
		return os.Stdout
	case "stderr":
		return os.Stderr
	default:
		file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return os.Stdout
		}
		return file
	}
}
