package main

import (
	"fmt"
	"os"
)

func getProfileOrSetDefault() string {
	return getValueOrSetDefault("PROFILE", "none-local")
}

func getDockerRegistryOrSetDefault() string {
	return getValueOrSetDefault("DOCKER_REGISTRY", "ghcr.io")
}

func getDockerRepositoryOrSetDefault() string {
	return getValueOrSetDefault("DOCKER_REPOSITORY", "valentinlutz")
}

func getProjectNameOrSetDefault() string {
	return getValueOrSetDefault("PROJECT_NAME", "order-service")
}

func getVersionOrSetDefault() string {
	return getValueOrSetDefault("VERSION", "latest")
}

func getFlywayUserOrSetDefault() string {
	return getValueOrSetDefault("FLYWAY_USER", "test")
}

func getFlywayPasswordOrSetDefault() string {
	return getValueOrSetDefault("FLYWAY_PASSWORD", "test")
}

func getValueOrSetDefault(key string, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		fmt.Printf("env '%s' not set, defaulting to '%s'\n", key, defaultValue)
		err := os.Setenv(key, defaultValue)
		if err != nil {
			panic(err)
		}
		return defaultValue
	}
	return value
}