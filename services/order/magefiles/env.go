package main

import (
	"fmt"
	"os"
)

func getProfileOrDefault() string {
	return getValueOrDefault("PROFILE", "none-local")
}

func getDockerRegistryOrDefault() string {
	return getValueOrDefault("DOCKER_REGISTRY", "ghcr.io")
}

func getDockerRepositoryOrDefault() string {
	return getValueOrDefault("DOCKER_REPOSITORY", "valentinlutz")
}

func getProjectNameOrDefault() string {
	return getValueOrDefault("PROJECT_NAME", "order-service")
}

func getVersionOrDefault() string {
	return getValueOrDefault("VERSION", "latest")
}

func getFlywayUserOrDefault() string {
	return getValueOrDefault("FLYWAY_USER", "test")
}

func getFlywayPasswordOrDefault() string {
	return getValueOrDefault("FLYWAY_PASSWORD", "test")
}

func getValueOrDefault(key string, defaultValue string) string {
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
