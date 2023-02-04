package config

import (
	"gopkg.in/yaml.v3"
	"monorepo/libraries/apputil/logging"
	"os"
)

type Config struct {
	Region      Region               `yaml:"region"`
	Environment Environment          `yaml:"environment"`
	Server      Server               `yaml:"server"`
	Database    Database             `yaml:"database"`
	Logger      logging.LoggerConfig `yaml:"logger"`
}

type Server struct {
	Port            int    `yaml:"port"`
	CertificatePath string `yaml:"certificate_path"`
	KeyPath         string `yaml:"key_path"`
}

type ClientConfig struct {
	Url string `yaml:"url"`
}

func New(path string) (*Config, error) {
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
