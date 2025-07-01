// Package main is the entry point for the Dynatrace bootstrapper.
package main

import (
	"os"

	bootstrapper "github.com/Dynatrace/dynatrace-bootstrapper/cmd"
	"github.com/spf13/afero"
)

func main() {
	err := bootstrapper.New(afero.NewOsFs()).Execute()
	if err != nil {
		os.Exit(1)
	}
}
