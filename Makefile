PHONY: build snapshot lint test
GIT_VERSION ?= $(shell git describe --tags --always --dirty="-dev")
DATE ?= $(shell date -u '+%Y-%m-%d %H:%M UTC')
VERSION_FLAGS := -s -w -X "main.version=$(GIT_VERSION)" -X "main.date=$(DATE)"
.DEFAULT_GOAL := build
DOCS_DIR := $(CURDIR)/documentation
GOTESTSUM_JUNITFILE ?= $(CURDIR)/-junit.xml
build:
	go build -trimpath -ldflags='$(VERSION_FLAGS) -extldflags -static' ./...

snapshot:
	goreleaser --snapshot --clean --skip=publish,sign

lint:
	golangci-lint run --config .golangci-lint.yml ./...

test:
	@gotestsum --format testname --junitfile junit.xml -- -coverprofile=cover.out ./...
