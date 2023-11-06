//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Docker mg.Namespace

// Builds the lambda container image | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
func (Docker) Build() error {
	dockerRegistry := getDockerRegistryOrSetDefault()
	dockerRepository := getDockerRepositoryOrSetDefault()
	projectName := getProjectNameOrSetDefault()
	version := getVersionOrSetDefault()

	return sh.RunV(
		"docker",
		"build",
		"--platform", "linux/amd64",
		"--file", "./app/Dockerfile",
		"--tag", dockerRegistry+"/"+dockerRepository+"/"+projectName+":"+version,
		"../",
	)
}

// Runs the lambda container image
func (Docker) Run() error {
	getDockerRegistryOrSetDefault()
	getDockerRepositoryOrSetDefault()
	getProjectNameOrSetDefault()
	getVersionOrSetDefault()

	os.Chdir("./deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./docker-compose.yaml",
		"up",
		"--force-recreate",
	)
}
