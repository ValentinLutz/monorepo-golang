//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Docker mg.Namespace

func (Docker) Up() error {
	getOrSetDefaultDatabaseEnvVars()

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

func (Docker) Down() error {
	getOrSetDefaultDatabaseEnvVars()

	os.Chdir("./deployment-docker")
	defer os.Chdir("..")

	return sh.RunV(
		"docker",
		"compose",
		"--file", "./docker-compose.yaml",
		"down",
	)
}
