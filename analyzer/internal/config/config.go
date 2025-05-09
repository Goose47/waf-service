// Package config provides config descriptions and function to read it.
package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
	"time"
)

var (
	errConfigNotFound = errors.New("config file does not exist")
	errBadConfig      = errors.New("can not read the config file")
)

// Config represents main application config.
type Config struct {
	Env       string          `yaml:"env"`
	Kafka     KafkaConfig     `yaml:"kafka"`
	Detection DetectionConfig `yaml:"detection"`
	Limiter   LimiterConfig   `yaml:"limiter"`
	GRPC      GRPCConfig      `yaml:"grpc"`
}

// GRPCConfig represents application server config.
type GRPCConfig struct {
	Port int `yaml:"port"`
}

// KafkaConfig represents application kafka config.
type KafkaConfig struct {
	Host           string `yaml:"host"`
	Topic          string `yaml:"topic"`
	DetectionTopic string `yaml:"detection_topic"`
	Port           int    `yaml:"port"`
}

// DetectionConfig represents detection service config.
type DetectionConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// LimiterConfig represents rate limiter config.
type LimiterConfig struct {
	MaxRequests int           `yaml:"max_requests"`
	Per         time.Duration `yaml:"per"`
}

// MustLoadPath loads config from configPath and panics on any errors.
func MustLoadPath(configPath string) *Config {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(errConfigNotFound)
	}

	var cfg Config

	err := cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		panic(errBadConfig)
	}

	return &cfg
}
