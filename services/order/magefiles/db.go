//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
)

type Db mg.Namespace

func (Db) Migrate() {
	flywayCredentials := getFlywayCredentials()
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
		"flyway/flyway:9.17.0-alpine",
		"clean",
		"migrate",
		"-user="+flywayCredentials.Username,
		"-password="+flywayCredentials.Password,
	)
}
