package config

import (
	"app/internal/database"
	"gopkg.in/yaml.v2"
	"os"
)

type Config struct {
	Server   ServerConfig    `yaml:"server"`
	Database database.Config `yaml:"database"`
}

type ServerConfig struct {
	Host    string        `yaml:"host"`
	Port    int           `yaml:"port"`
	Timeout TimeoutConfig `yaml:"timeout"`
}

type TimeoutConfig struct {
	Read  int `yaml:"read"`
	Write int `yaml:"write"`
	Idle  int `yaml:"idle"`
}

func NewConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := &Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
