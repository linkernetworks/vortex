# go tool paths
GO_COVER := $(shell which cover)
GO_COBERTURA := $(shell which gocover-cobertura)
GO_VENDOR := $(shell which govendor)

GO_FILES := $(shell find src -type f -iname "*.go")

CONFIG_FILE = config/k8s.json

SHELL := /bin/bash

GCP_DOCKER_REGISTRY = asia.gcr.io

all:

test:

deps:
	go get ./...

vortex:
	go build ./src/cmd/vortex

run: vortex
	./vortex -config config/local.json -port 7890

image:
	docker build -f dockerfiles/Dockerfile .
