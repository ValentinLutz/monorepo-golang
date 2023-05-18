package config

import (
	"fmt"
)

type Environment string

const (
	LOCAL     Environment = "LOCAL"
	CONTAINER Environment = "CONTAINER"
	TEST      Environment = "TEST"
	PROD      Environment = "PROD"
)

// Implements the yaml.Unmarshaler interface.
func (environment *Environment) UnmarshalYAML(unmarshal func(any) error) error {
	var environmentString string
	err := unmarshal(&environmentString)
	if err != nil {
		return err
	}

	unmarshaledEnvironment := Environment(environmentString)

	switch unmarshaledEnvironment {
	case LOCAL, CONTAINER, TEST, PROD:
		*environment = unmarshaledEnvironment
		return nil
	}
	return fmt.Errorf("environment '%v' is invalid", unmarshaledEnvironment)
}
