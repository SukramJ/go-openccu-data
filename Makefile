# SPDX-License-Identifier: MIT
# go-openccu-data — developer Makefile
#
# Tabs are required by GNU make. The whitespace rules below pin sane
# shell behaviour so a failing recipe step actually aborts the target
# instead of silently moving on.

SHELL := /usr/bin/env bash
.SHELLFLAGS := -euo pipefail -c
.DEFAULT_GOAL := help

GO      ?= go
GOFUMPT ?= gofumpt

export CGO_ENABLED := 0

# Path to a local openccu-data checkout for snapshot regeneration.
OPENCCU_DATA ?= ../openccu-data

.PHONY: help
help: ## show this help
	@awk 'BEGIN {FS = ":.*## "} /^[a-zA-Z_-]+:.*## / {printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: setup
setup: ## install developer tooling (gofumpt)
	$(GO) install mvdan.cc/gofumpt@latest

.PHONY: test
test: ## run the full test suite with race detector
	CGO_ENABLED=1 $(GO) test -race -count=1 -timeout=120s ./...

.PHONY: test-cover
test-cover: ## run tests + coverage report
	CGO_ENABLED=1 $(GO) test -race -count=1 -covermode=atomic -coverprofile=coverage.out ./...
	$(GO) tool cover -func=coverage.out | tail -20

.PHONY: vet
vet: ## run go vet
	$(GO) vet ./...

.PHONY: fmt
fmt: ## format with gofumpt (writes in place)
	$(GOFUMPT) -w .

.PHONY: fmt-check
fmt-check: ## fail when sources are not gofumpt-clean
	@diff=$$($(GOFUMPT) -l .); \
	if [ -n "$$diff" ]; then \
	  echo "gofumpt would rewrite:"; echo "$$diff"; exit 1; \
	fi

.PHONY: tidy
tidy: ## sync go.mod
	$(GO) mod tidy

.PHONY: check
check: vet fmt-check test ## the pre-commit / pre-push gate

.PHONY: regenerate
regenerate: ## regenerate data/ from $(OPENCCU_DATA) and re-run the tests
	script/regenerate.sh $(OPENCCU_DATA)
	$(MAKE) --no-print-directory test

.PHONY: snapshot-version
snapshot-version: ## print the embedded openccu-data snapshot version
	@sed -n 's/^const SnapshotVersion = "\(.*\)"$$/\1/p' openccudata.go

.PHONY: clean
clean: ## remove build artefacts
	rm -f coverage.out
