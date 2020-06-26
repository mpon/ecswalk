package main

import "github.com/mpon/ecswalk/internal/command"

// Version is dynamically set by the toolchain or overridden by the Makefile.
var Version = "DEV"

func main() {
	command.Execute(Version)
}
