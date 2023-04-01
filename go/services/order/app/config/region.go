package config

import (
	"fmt"
)

type Region string

const (
	NONE Region = "NONE"
	EU   Region = "EU"
	US   Region = "US"
)

// Implements the yaml.Unmarshaler interface.
func (r *Region) UnmarshalYAML(unmarshal func(any) error) error {
	var regionString string
	err := unmarshal(&regionString)
	if err != nil {
		return err
	}

	region := Region(regionString)

	switch region {
	case NONE, EU, US:
		*r = region
		return nil
	}
	return fmt.Errorf("region '%v' is invalid", region)
}
