//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"os"
)

type App mg.Namespace

func (App) Run() error {
	getVersionOrSetDefault()

	mg.Deps(Dep.Copy)

	os.Chdir("./app")
	defer os.Chdir("..")

	return sh.RunV("go", "run", "main.go")
}

func (App) Build() error {
	os.Chdir("./app")
	defer os.Chdir("..")

	return sh.RunV("go", "build")
}
