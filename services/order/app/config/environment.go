package config

import "fmt"

type Environment string

const (
	DEV  Environment = "DEV"
	TEST Environment = "TEST"
	PROD Environment = "PROD"
)

func (env *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var environmentString string
	err := unmarshal(&environmentString)
	if err != nil {
		return err
	}

	environment := Environment(environmentString)

	switch environment {
	case DEV, TEST, PROD:
		*env = environment
		return nil
	}
	return fmt.Errorf("environment '%v' is invalid", environment)
}
