package config

import "fmt"

type Environment string

const (
	LOCAL     Environment = "LOCAL"
	CONTAINER Environment = "CONTAINER"
	TEST      Environment = "TEST"
	PROD      Environment = "PROD"
)

func (e *Environment) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var environmentString string
	err := unmarshal(&environmentString)
	if err != nil {
		return err
	}

	environment := Environment(environmentString)

	switch environment {
	case LOCAL, CONTAINER, TEST, PROD:
		*e = environment
		return nil
	}
	return fmt.Errorf("environment '%v' is invalid", environment)
}
