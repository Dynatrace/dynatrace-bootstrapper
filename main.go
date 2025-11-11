package main

import (
	"os"

	bootstrapper "github.com/Dynatrace/dynatrace-bootstrapper/cmd"
)

func main() {
	err := bootstrapper.New().Execute()
	if err != nil {
		os.Exit(1)
	}
}
