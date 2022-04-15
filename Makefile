include .make/help.mk
include .make/docker.mk
include .make/database.mk
include .make/app.mk

PROJECT_NAME ?= golang-reference-project
PROFILE ?= none-dev
FLYWAY_USER ?= test
FLYWAY_PASSWORD ?= test