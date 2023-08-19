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

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (region *Region) UnmarshalYAML(unmarshal func(any) error) error {
	var regionString string
	err := unmarshal(&regionString)
	if err != nil {
		return err
	}

	unmarshalledRegion := Region(regionString)

	switch unmarshalledRegion {
	case NONE, EU, US:
		*region = unmarshalledRegion
		return nil
	}
	return fmt.Errorf("region '%v' is invalid", unmarshalledRegion)
}
