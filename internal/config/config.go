package config

import (
	"time"

	"github.com/spf13/viper"
)

// Config holds application configuration
type Config struct {
	Group     string        `mapstructure:"group"`
	Interface string        `mapstructure:"interface"`
	Interval  time.Duration `mapstructure:"interval"`
}

// Load reads configuration from file and environment
func Load() (*Config, error) {
	var cfg Config

	// Set defaults
	viper.SetDefault("group", "224.1.1.1:9999")
	viper.SetDefault("interval", time.Second)

	// Environment variables
	viper.SetEnvPrefix("MULTICAST")
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
