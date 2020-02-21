.PHONY: all build build-local unit-test deps lint check-goos
SHELL=/bin/bash

version=$(shell cat version/version.go | grep Version | cut -d'"' -f2)

ifeq ($(OS),Windows_NT)
    OSNAME = windows
else
    UNAME_S := $(shell uname -s)
    ifeq ($(UNAME_S),Linux)
        OSNAME = linux
    endif
    ifeq ($(UNAME_S),Darwin)
        OSNAME = darwin
    endif
endif

ifdef os
  OSNAME=$(os)
endif

all: unit-test lint build

build-local: deps unit-test
	operator-sdk build --go-build-args '-ldflags=-s -ldflags=-w' xo.90poe.io/elasticsearch-operator:$(version)

build:
	CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -mod=vendor -ldflags="-s -w" -a -o ./artifacts/manager ./cmd/manager
	mv ./artifacts/manager ./artifacts/manager-unpacked
	upx -q -o ./artifacts/manager ./artifacts/manager-unpacked
	rm -rf ./artifacts/manager-unpacked

deps:
	operator-sdk generate crds
	operator-sdk generate k8s
	go mod vendor

unit-test:
	go test -v -parallel=2 -mod=vendor -cover $$(go list ./...)

lint:
	# (cd /tmp/; go get -u github.com/golangci/golangci-lint/...)
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run
