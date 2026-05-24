package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadUsesConfigFileAndEnvironmentOverrides(t *testing.T) {
	dir := t.TempDir()
	configPath := filepath.Join(dir, "config.yaml")
	configContent := []byte("environment: staging\napi_server:\n  port: 8085\nlogging:\n  level: info\n")
	if err := os.WriteFile(configPath, configContent, 0o600); err != nil {
		t.Fatalf("write config: %v", err)
	}

	t.Setenv(envConfigPath, configPath)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.Environment != "staging" {
		t.Fatalf("expected environment staging, got %q", cfg.Environment)
	}
	if cfg.APIServer.Port != 8085 {
		t.Fatalf("expected api server port from file 8085, got %d", cfg.APIServer.Port)
	}
	if cfg.Logging.Level != "info" {
		t.Fatalf("expected logging level from file info, got %q", cfg.Logging.Level)
	}
}

func TestLoadReturnsDefaultsWhenConfigFileIsMissing(t *testing.T) {
	t.Setenv(envConfigPath, "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() returned error: %v", err)
	}

	if cfg.APIServer.Port != DefaultConfig().APIServer.Port {
		t.Fatalf("expected default api server port %d, got %d", DefaultConfig().APIServer.Port, cfg.APIServer.Port)
	}
	if cfg.Environment != DefaultConfig().Environment {
		t.Fatalf("expected default environment %q, got %q", DefaultConfig().Environment, cfg.Environment)
	}
}
