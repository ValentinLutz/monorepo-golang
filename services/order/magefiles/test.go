//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Test mg.Namespace

func (Test) Unit() {
	mg.Deps(Dep.Copy, Dep.Generate)

	os.Chdir("./app")
	defer os.Chdir("..")

	sh.RunV("go", "test", "./...")
}

func (Test) Lint() {
	mg.Deps(Dep.Copy, Dep.Generate)

	os.Chdir("./app")
	defer os.Chdir("..")

	sh.RunV("golangci-lint", "run")
}

func (Test) Spectral() {
	os.Chdir("./api-definition")
	defer os.Chdir("..")

	sh.RunV(
		"docker",
		"run",
		"--rm",
		"-it",
		"--volume", "./spectral.ruleset.yaml:/tmp/ruleset.yaml",
		"--volume", "./order_api.yaml:/tmp/order_api.yaml",
		"stoplight/spectral:6.7.0",
		"--ruleset", "/tmp/ruleset.yaml",
		"lint", "/tmp/order_api.yaml",
	)
}

func (Test) Smoke() {
	getProfileOrSetDefault()

	mg.Deps(Dep.Generate)

	os.Chdir("./test-smoke")
	defer os.Chdir("..")

	sh.RunV("go", "test", "-count=1", "./...")
}

func (Test) Functional() {
	getProfileOrSetDefault()

	mg.Deps(Dep.Generate)

	os.Chdir("./test-integration")
	defer os.Chdir("..")

	sh.RunV("go", "test", "-count=1", "./...")
}

func (Test) Coverage() {
	getProfileOrSetDefault()
	os.RemoveAll("./test-integration/coverage")

	mg.Deps(Docker.Up)
	mg.Deps(Test.Functional)
	mg.Deps(Docker.Down)

	os.Chdir("./test-integration")
	defer os.Chdir("..")

	sh.RunV("go", "tool", "covdata", "percent", "-i", "./coverage")
}

func (Test) Coveragehtml() {
	os.Chdir("./test-integration")
	defer os.Chdir("..")

	sh.RunV("go", "tool", "covdata", "textfmt", "-i", "./coverage", "-o", "./coverage/coverage.txt")
	sh.RunV("go", "tool", "cover", "-html=./coverage/coverage.txt", "-o", "./coverage/coverage.html")
}

func (Test) Load() {
	os.Chdir("./test-load")
	defer os.Chdir("..")

	sh.RunV(
		"docker",
		"run",
		"-it",
		"--rm",
		"--network", "host",
		"--volume", "./script.js:/k6/script.js",
		"grafana/k6:0.39.0",
		"run", "/k6/script.js",
	)
}
