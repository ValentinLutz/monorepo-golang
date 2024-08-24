//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Dep mg.Namespace

func (Dep) Copy() error {
	err := sh.RunV("ginstall", "-D", "./config/app.config.none-dev.yaml", "./app/config/config.yaml")
	if err != nil {
		return err
	}
	err = sh.RunV("ginstall", "-D", "./config/app.private.none-dev.key", "./app/config/app.private.none-dev.key")
	if err != nil {
		return err
	}
	return sh.RunV("ginstall", "-D", "./config/app.public.none-dev.crt", "./app/config/app.public.none-dev.crt")
}

func (Dep) Generate() error {
	return sh.RunV("go", "generate", "./...")
}
