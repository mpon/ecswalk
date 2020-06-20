# The binary to build (just the basename).
BIN ?= ecswalk

# This repo's root import path (under GOPATH).
PKG := github.com/mpon/ecswalk

VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS := -X $(PKG)/internal/command.Version=$(VERSION)

build:
	go build \
		-ldflags "$(LDFLAGS)" \
		-o $(BIN) \
		$(PKG)/cmd/$(BIN)

test:
	go test -v ./...
