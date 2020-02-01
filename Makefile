.PHONY: all build deps lint check-goos
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

all: unit_test lint build

build: unit_test
	operator-sdk build --go-build-args '-ldflags=-s -ldflags=-w' xo.90poe.io/elasticsearch-operator:$(version)

deps:
	go mod vendor

unit_test:
	go test -v -mod=vendor -cover $$(go list ./...)

lint:
	# (cd /tmp/; go get -u github.com/golangci/golangci-lint/...)
	# go install github.com/golangci/golangci-lint/cmd/golangci-lint
	golangci-lint run
