package configs

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
)

type Config struct {
	Logger   LoggerConfig   `yaml:"logger"`
	Storage  StorageConfig  `yaml:"storage"`
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

func NewConfig(configFile io.Reader) (*Config, error) {
	config := &Config{}

	decoder := yaml.NewDecoder(configFile)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return config, nil
}
