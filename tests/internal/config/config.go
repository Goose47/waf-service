package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

type Config struct {
	Analyzer  GRPCConfig    `yaml:"analyzer"`
	Detection GRPCConfig    `yaml:"detection"`
	WAF       GRPCConfig    `yaml:"waf"`
	Limiter   LimiterConfig `yaml:"limiter"`
}

type GRPCConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// LimiterConfig represents rate limiter config.
type LimiterConfig struct {
	MaxRequests int           `yaml:"max_requests"`
	Per         time.Duration `yaml:"per"`
}

// MustLoadPath loads config from configPath and panics on any errors
func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Errorf("config file does not exist: %s", configPath))
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(fmt.Errorf("failed to read config from %s: %s", configPath, err.Error()))
	}

	return &cfg
}
