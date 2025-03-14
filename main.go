package main

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd"
	"github.com/spf13/afero"
)

func main() {
	err := cmd.New(afero.NewOsFs()).Execute()
	if err != nil {
		os.Exit(1)
	}
}
