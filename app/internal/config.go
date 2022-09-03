package internal

import (
	"app/external/database"
	"app/internal/config"
	"app/internal/util"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Region      config.Region      `yaml:"region"`
	Environment config.Environment `yaml:"environment"`
	Server      ServerConfig       `yaml:"server"`
	Logger      util.LoggerConfig  `yaml:"logger"`
	Database    database.Config    `yaml:"database"`
	Client      Client             `yaml:"client"`
}

type ServerConfig struct {
	Port    int           `yaml:"port"`
	Timeout TimeoutConfig `yaml:"timeout"`
}

type TimeoutConfig struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
	Idle  int `yaml:"idle"`
}

type Client struct {
	PaymentClient ClientConfig `yaml:"payment"`
}

type ClientConfig struct {
	Url string `yaml:"url"`
}

func NewConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return &Config{}, err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(file)

	decoder := yaml.NewDecoder(file)
	var decodedConfig *Config
	err = decoder.Decode(&decodedConfig)
	if err != nil {
		return &Config{}, err
	}

	return decodedConfig, nil
}
