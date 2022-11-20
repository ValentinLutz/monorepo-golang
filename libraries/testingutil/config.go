package testingutil

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	BaseURL  string `yaml:"base_url"`
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Database string `yaml:"database"`
	Username string `yaml:"user"`
	Password string `yaml:"password"`
}

func ParseFile[T any](path string) (*T, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(file)

	decoder := yaml.NewDecoder(file)
	var decodedConfig *T
	err = decoder.Decode(&decodedConfig)
	if err != nil {
		return nil, err
	}

	return decodedConfig, nil
}
