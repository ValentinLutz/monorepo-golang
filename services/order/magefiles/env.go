package main

import (
	"fmt"
	"os"
)

func getProfileOrDefault() string {
	key := "PROFILE"
	defaultValue := "none-local"

	profile, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("env '%s' not set, defaulting to '%s'\n", key, defaultValue)
		err := os.Setenv(key, defaultValue)
		if err != nil {
			panic(err)
		}
		return defaultValue
	}
	return profile
}

type FlywayCredentials struct {
	Username string
	Password string
}

func getFlywayCredentials() FlywayCredentials {
	return FlywayCredentials{
		Username: "test",
		Password: "test",
	}
}
