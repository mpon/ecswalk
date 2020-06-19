# The binary to build (just the basename).
BIN ?= ecswalk

# This repo's root import path (under GOPATH).
PKG := github.com/mpon/ecswalk

build:
	go build \
		-o ${BIN} \
		${PKG}/cmd/${BIN}

test:
	go test -v ./...
