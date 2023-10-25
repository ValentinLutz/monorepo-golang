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
		"--file", "./app/app.dockerfile",
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
		"--detach",
		"--force-recreate",
		"--wait",
	)
}

func (Docker) Appup() error {
	mg.Deps(Dep.Generate)

	os.Chdir("./deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./app.docker-compose.yaml",
		"up",
		"--force-recreate",
		"--build",
	)
}

func (Docker) Testup() error {
	mg.Deps(Dep.Generate)

	os.Chdir("./deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./test.docker-compose.yaml",
		"up",
		"--force-recreate",
		"--build",
		"-d",
	)
}

func (Docker) Testdown() error {
	os.Chdir("./deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./test.docker-compose.yaml",
		"down",
	)
}

func (Docker) Down() error {
	os.Chdir("./deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./docker-compose.yaml",
		"--file", "./app.docker-compose.yaml",
		"down",
	)
}
