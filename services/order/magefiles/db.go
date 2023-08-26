//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
)

const flywayDockerImage = "flyway/flyway:9.21.2-alpine"

type Db mg.Namespace

// Cleans flyway schema from the database | PROFILE, FLYWAY_USER, FLYWAY_PASSWORD
func (Db) Clean() {
	flywayUser := getFlywayUserOrDefault()
	flywayPassword := getFlywayPasswordOrDefault()
	profile := getProfileOrDefault()

	os.Chdir("./migration-database")
	defer os.Chdir("..")

	sh.RunV(
		"docker",
		"run",
		"-it",
		"--rm",
		"--network", "host",
		"--volume", "./migration:/flyway/sql/migration",
		"--volume", "./"+profile+".conf:/flyway/conf/flyway.conf",
		flywayDockerImage,
		"clean",
		"-user="+flywayUser,
		"-password="+flywayPassword,
	)
}

// Migrates flyway schema to the database | PROFILE, FLYWAY_USER, FLYWAY_PASSWORD
func (Db) Migrate() {
	flywayUser := getFlywayUserOrDefault()
	flywayPassword := getFlywayPasswordOrDefault()
	profile := getProfileOrDefault()

	os.Chdir("./migration-database")
	defer os.Chdir("..")
	sh.RunV(
		"docker",
		"run",
		"-it",
		"--rm",
		"--network", "host",
		"--volume", "./migration:/flyway/sql/migration",
		"--volume", "./"+profile+".conf:/flyway/conf/flyway.conf",
		flywayDockerImage,
		"migrate",
		"-user="+flywayUser,
		"-password="+flywayPassword,
	)
}
