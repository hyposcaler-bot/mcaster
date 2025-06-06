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
	TTL       int           `mapstructure:"ttl"`
	SPort     int           `mapstructure:"sport"`
	DPort     int           `mapstructure:"dport"`
}

// Load reads configuration from file and environment
func Load() (*Config, error) {
	var cfg Config

	// Set defaults
	viper.SetDefault("group", "239.23.23.23:2323")
	viper.SetDefault("interval", time.Second)
	viper.SetDefault("ttl", 1)
	viper.SetDefault("sport", 0)
	viper.SetDefault("dport", 0)

	// Environment variables
	viper.SetEnvPrefix("MULTICAST")
	viper.AutomaticEnv()

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
