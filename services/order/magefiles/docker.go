//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
)

type Docker mg.Namespace

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
