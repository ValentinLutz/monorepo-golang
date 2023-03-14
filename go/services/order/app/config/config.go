package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"monorepo/libraries/apputil/infastructure"
	"monorepo/libraries/apputil/logging"
	"os"
)

type Config struct {
	Version     string
	ServiceName string
	Region      Region                       `yaml:"region"`
	Environment Environment                  `yaml:"environment"`
	Server      infastructure.ServerConfig   `yaml:"server"`
	Database    infastructure.DatabaseConfig `yaml:"database"`
	Logger      logging.LoggerConfig         `yaml:"logger"`
}

func New(path string) (*Config, error) {
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
	var decodedConfig *Config
	err = decoder.Decode(&decodedConfig)
	if err != nil {
		return nil, err
	}

	version, ok := os.LookupEnv("VERSION")
	if !ok {
		return nil, errors.New("failed to load the environment variable 'VERSION'")
	}
	decodedConfig.Version = version

	projectName, ok := os.LookupEnv("PROJECT_NAME")
	if !ok {
		return nil, errors.New("failed to load the environment variable 'PROJECT_NAME'")
	}
	decodedConfig.ServiceName = projectName

	return decodedConfig, nil
}
