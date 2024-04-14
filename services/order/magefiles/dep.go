//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Dep mg.Namespace

func (Dep) Copy() error {
	err := sh.RunV("ginstall", "-D", "./config/app.config.none-local.yaml", "./app/config/config.yaml")
	if err != nil {
		return err
	}
	err = sh.RunV("ginstall", "-D", "./config/app.key", "./app/config/app.key")
	if err != nil {
		return err
	}
	return sh.RunV("ginstall", "-D", "./config/app.crt", "./app/config/app.crt")
}

func (Dep) Generate() error {
	err := sh.RunV("oapi-codegen", "--config", "./api-definition/app.model.yaml", "./api-definition/order_api.yaml")
	if err != nil {
		return err
	}
	err = sh.RunV("oapi-codegen", "--config", "./api-definition/app.server.yaml", "./api-definition/order_api.yaml")
	if err != nil {
		return err
	}
	return sh.RunV("oapi-codegen", "--config", "./api-definition/test.client.yaml", "./api-definition/order_api.yaml")
}
