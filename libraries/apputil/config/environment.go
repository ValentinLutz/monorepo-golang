package config

import (
	"fmt"
)

type Environment string

const (
	DEV  Environment = "DEV"
	TEST Environment = "TEST"
	PROD Environment = "PROD"
)

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (environment *Environment) UnmarshalYAML(unmarshal func(any) error) error {
	var environmentString string
	err := unmarshal(&environmentString)
	if err != nil {
		return err
	}

	unmarshalledEnvironment := Environment(environmentString)

	switch unmarshalledEnvironment {
	case DEV, TEST, PROD:
		*environment = unmarshalledEnvironment
		return nil
	}
	return fmt.Errorf("environment '%v' is invalid", unmarshalledEnvironment)
}
