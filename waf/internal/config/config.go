// Package config provides config descriptions and function to read it.
package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
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
	Analyzer  AnalyzerConfig  `yaml:"analyzer"`
	GRPC      GRPCConfig      `yaml:"grpc"`
}

// GRPCConfig represents application server config.
type GRPCConfig struct {
	Port int `yaml:"port"`
}

// KafkaConfig represents application kafka config.
type KafkaConfig struct {
	Host          string `yaml:"host"`
	AnalyzerTopic string `yaml:"analyzer_topic"`
	Port          int    `yaml:"port"`
}

// DetectionConfig represents detection gRPC service config.
type DetectionConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

// AnalyzerConfig represents analyzer gRPC service config.
type AnalyzerConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
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
