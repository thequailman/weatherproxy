BINDIR = $(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))/.bin
GO_VERSION ?= 1.14.2
GOLANGCI_LINT_VERSION ?= 1.26.0
SHELL := /bin/bash

export GOROOT = $(BINDIR)/go/lib
export GOPATH = $(BINDIR)/go
export PATH := $(BINDIR)/go/bin:$(BINDIR)/go/lib/bin:$(BINDIR):$(PATH)

.PHONY: go golangci-lint node

go:
	@if ! ${BINDIR}/go/lib/bin/go version | grep $(GO_VERSION) > /dev/null; then \
		rm -rf .bin/go/lib; \
		mkdir -p .bin/go/lib; \
		curl -s -L https://dl.google.com/go/go$(GO_VERSION).linux-amd64.tar.gz | tar -C .bin/go/lib --strip-components=1 -xz; \
		go get golang.org/x/tools/gopls; \
	fi;

golangci-lint: go
	@if ! ${BINDIR}/golangci-lint version 2>&1 | grep $(GOLANGCI_LINT_VERSION) > /dev/null; then \
		mkdir -p .bin; \
		curl -s -L https://github.com/golangci/golangci-lint/releases/download/v$(GOLANGCI_LINT_VERSION)/golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64.tar.gz | tar --no-same-owner -C .bin -xz --strip-components=1 golangci-lint-$(GOLANGCI_LINT_VERSION)-linux-amd64/golangci-lint; \
		chmod +x .bin/golangci-lint; \
	fi;
