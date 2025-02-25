package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"os"
)

type Config struct {
	Env   string      `yaml:"env"`
	GRPC  GRPCConfig  `yaml:"grpc"`
	Kafka KafkaConfig `yaml:"kafka"`
}

type GRPCConfig struct {
	Port int `yaml:"port"`
}

type KafkaConfig struct {
	Host           string `yaml:"host"`
	Port           int    `yaml:"port"`
	Topic          string `yaml:"topic"`
	DetectionTopic string `yaml:"detection_topic"`
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
