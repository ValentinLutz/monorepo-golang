//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Lambda mg.Namespace

func (Lambda) Invoke() error {
	return sh.RunV("sam", "local", "invoke", "--docker-network", "facts-lambda-network")
}

func (Lambda) Build() error {
	return sh.RunV("sam", "build")
}
