//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Docker mg.Namespace

// Builds the container image | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
func (Docker) Build() error {
	dockerRegistry := getDockerRegistryOrSetDefault()
	dockerRepository := getDockerRepositoryOrSetDefault()
	projectName := getProjectNameOrSetDefault()
	version := getVersionOrSetDefault()

	mg.Deps(Dep.Generate)

	return sh.RunV(
		"docker",
		"build",
		"--file", "./app/Dockerfile",
		"--tag", dockerRegistry+"/"+dockerRepository+"/"+projectName+":"+version,
		"../../",
	)
}

// Publish the container image | DOCKER_REGISTRY, DOCKER_REPOSITORY, PROJECT_NAME, VERSION
func (Docker) Push() error {
	dockerRegistry := getDockerRegistryOrSetDefault()
	dockerRepository := getDockerRepositoryOrSetDefault()
	projectName := getProjectNameOrSetDefault()
	version := getVersionOrSetDefault()

	return sh.RunV(
		"docker",
		"push",
		dockerRegistry+"/"+dockerRepository+"/"+projectName+":"+version,
	)
}

func (Docker) Up() error {
	os.Chdir("./deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./docker-compose.yaml",
		"up",
		"--force-recreate",
		"--build",
		"--wait",
		//"--detach",
	)
}

func (Docker) Down() error {
	os.Chdir("./deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./docker-compose.yaml",
		"down",
	)
}
