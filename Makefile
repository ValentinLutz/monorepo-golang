include .make/help.mk
include .make/docker.mk
include .make/database.mk
include .make/app.mk
include .make/test.mk

PROJECT_NAME ?= golang-reference-project
PROFILE ?= none-dev
FLYWAY_USER ?= test
FLYWAY_PASSWORD ?= test
VERSION ?= 0.1.2
DOCKER_REGISTRY ?= ghcr.io
DOCKER_REPOSITORY ?= valentinlutz