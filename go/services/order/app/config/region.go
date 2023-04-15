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
func (region *Region) UnmarshalYAML(unmarshal func(any) error) error {
	var regionString string
	err := unmarshal(&regionString)
	if err != nil {
		return err
	}

	unmarshaledRegion := Region(regionString)

	switch unmarshaledRegion {
	case NONE, EU, US:
		*region = unmarshaledRegion
		return nil
	}
	return fmt.Errorf("region '%v' is invalid", unmarshaledRegion)
}
