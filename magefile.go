//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Dep mg.Namespace

// Installs all dependencies for the project
func (Dep) Install() error {
	err := sh.RunV("go", "install", "github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen@v2.1.0")
	if err != nil {
		return err
	}
	return sh.RunV("go", "install", "github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2")
}
