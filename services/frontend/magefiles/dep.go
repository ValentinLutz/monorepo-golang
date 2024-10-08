//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Dep mg.Namespace

func (Dep) Copy() error {
	err := sh.RunV("install", "-D", "./config/app.config.none-dev.yaml", "./app/config/config.yaml")
	if err != nil {
		return err
	}
	err = sh.RunV("install", "-D", "./config/app.key", "./app/config/app.key")
	if err != nil {
		return err
	}
	return sh.RunV("install", "-D", "./config/app.crt", "./app/config/app.crt")
}
