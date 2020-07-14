# The binary to build (just the basename).
BIN ?= ecswalk

VERSION ?= $(shell git rev-parse --short HEAD)
LDFLAGS := -X main.Version=$(VERSION)

build:
	go build -ldflags "$(LDFLAGS)" -o $(BIN) .
test:
	go test -v ./...
lint:
	golangci-lint run
