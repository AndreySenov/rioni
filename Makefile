# Copyright (c) 2026 Andrey Senov
# SPDX-License-Identifier: Apache-2.0

GO ?= go
GIT ?= git
GOLICENSES ?= go-licenses
GORELEASER ?= goreleaser
BINARY_NAME := rioni
BIN_PATH := bin
CERT_PATH := cert
DIST_PATH := dist
CMD_PATH := ./cmd/$(BINARY_NAME)

.PHONY: clean
clean:
	rm -rf $(BIN_PATH)
	rm -rf $(CERT_PATH)
	rm -rf $(DIST_PATH)
	$(GO) clean

.PHONY: tidy
tidy:
	$(GO) mod tidy
	$(GIT) diff --exit-code -- go.mod go.sum

.PHONY: vet
vet:
	$(GO) vet ./...

.PHONY: test
test:
	$(GO) test ./...

.PHONY: check
check: tidy vet test

.PHONY: build
build:
	$(GO) build -o $(BIN_PATH)/$(BINARY_NAME) $(CMD_PATH)

.PHONY: run
run:
	@set -a; \
    [ -f .env ] && source ./.env; \
    set +a; \
	$(GO) run $(CMD_PATH) $(ARGS)

.PHONY: report-licenses
report-licenses:
	$(GOLICENSES) report $(CMD_PATH)

.PHONY: build-release-snapshot
build-release-snapshot:
	$(GORELEASER) check
	$(GORELEASER) release --snapshot --clean
