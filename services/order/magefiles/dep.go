//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Dep mg.Namespace

func (Dep) Install() error {
	err := sh.RunV("go", "install", "github.com/deepmap/oapi-codegen/cmd/oapi-codegen@v1.12.4")
	if err != nil {
		return err
	}
	return sh.RunV("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@v1.54.2")
}

func (Dep) Copy() error {
	err := sh.RunV("install", "-D", "./config/app.config.none-local.yaml", "./app/config/config.yaml")
	if err != nil {
		return err
	}
	err = sh.RunV("install", "-D", "./config/app.key", "./app/config/app.key")
	if err != nil {
		return err
	}
	return sh.RunV("install", "-D", "./config/app.crt", "./app/config/app.crt")
}

func (Dep) Generate() error {
	mg.Deps(Dep.Install)

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
