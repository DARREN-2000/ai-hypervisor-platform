package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/viper"
)

const (
	defaultConfigName = "config"
	defaultConfigType = "yaml"
	envPrefix         = "AIHYPERVISOR"
	envConfigPath     = "AIHYPERVISOR_CONFIG_PATH"
)

// Load reads configuration from file and environment variables.
func Load() (*PlatformConfig, error) {
	cfg := DefaultConfig()
	loader := viper.New()

	if configPath := os.Getenv(envConfigPath); configPath != "" {
		loader.SetConfigFile(configPath)
	} else {
		loader.SetConfigName(defaultConfigName)
		loader.SetConfigType(defaultConfigType)
		loader.AddConfigPath("/etc/aihypervisor")
		loader.AddConfigPath("$HOME/.aihypervisor")
		loader.AddConfigPath(".")
	}

	replacer := strings.NewReplacer(".", "_")
	loader.SetEnvPrefix(envPrefix)
	loader.SetEnvKeyReplacer(replacer)
	loader.AutomaticEnv()

	if err := loader.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	if err := loader.Unmarshal(cfg, func(dc *mapstructure.DecoderConfig) {
		dc.TagName = "yaml"
	}); err != nil {
		return nil, fmt.Errorf("unmarshal config: %w", err)
	}

	return cfg, nil
}
