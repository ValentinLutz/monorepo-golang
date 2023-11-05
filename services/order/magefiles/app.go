//go:build mage
// +build mage

package main

import (
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type App mg.Namespace

func (App) Run() error {
	getVersionOrSetDefault()

	mg.Deps(Dep.Copy, Dep.Generate)

	os.Chdir("./app")
	defer os.Chdir("..")

	return sh.RunV("go", "run", "main.go")
}

func (App) Build() error {
	mg.Deps(Dep.Generate)

	os.Chdir("./app")
	defer os.Chdir("..")

	return sh.RunV("go", "build")
}
