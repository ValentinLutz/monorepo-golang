package main

import (
	"fmt"
	"os"
)

func getDockerRegistryOrSetDefault() string {
	return getValueOrSetDefault("DOCKER_REGISTRY", "ghcr.io")
}

func getDockerRepositoryOrSetDefault() string {
	return getValueOrSetDefault("DOCKER_REPOSITORY", "valentinlutz")
}

func getProjectNameOrSetDefault() string {
	return getValueOrSetDefault("PROJECT_NAME", "test-lambda")
}

func getVersionOrSetDefault() string {
	return getValueOrSetDefault("VERSION", "latest")
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
