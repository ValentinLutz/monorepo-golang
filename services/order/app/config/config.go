package config

import (
	"errors"
	"monorepo/libraries/apputil/config"
	"monorepo/libraries/apputil/infastructure"
	"monorepo/libraries/apputil/logging"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version     string
	Region      config.Region                `yaml:"region"`
	Environment config.Environment           `yaml:"environment"`
	Server      infastructure.ServerConfig   `yaml:"server"`
	Database    infastructure.DatabaseConfig `yaml:"database"`
	Logger      logging.LoggerConfig         `yaml:"logger"`
}

func New(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer func(file *os.File) {
		closeErr := file.Close()
		if closeErr != nil {
			err = closeErr
		}
	}(file)

	decoder := yaml.NewDecoder(file)
	var decodedConfig Config
	err = decoder.Decode(&decodedConfig)
	if err != nil {
		return Config{}, err
	}

	version, ok := os.LookupEnv("VERSION")
	if !ok {
		return Config{}, errors.New("failed to load env 'VERSION'")
	}
	decodedConfig.Version = version

	return decodedConfig, nil
}
