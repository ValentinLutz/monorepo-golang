package config

import (
	"errors"
	"gopkg.in/yaml.v3"
	"monorepo/libraries/apputil/logging"
	"os"
)

type Config struct {
	Version     string
	ServiceName string
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
