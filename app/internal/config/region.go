package config

import "fmt"

type Region string

const (
	NONE Region = "NONE"
	EU          = "EU"
)

func (r *Region) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var regionString string
	err := unmarshal(&regionString)
	if err != nil {
		return err
	}

	region := Region(regionString)

	switch region {
	case NONE, EU:
		*r = region
		return nil
	}
	return fmt.Errorf("region is invalid: %s", region)
}
