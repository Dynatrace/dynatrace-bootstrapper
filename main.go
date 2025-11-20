package main

import (
	"os"

	bootstrapper "github.com/Dynatrace/dynatrace-bootstrapper/cmd/init"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/serverless"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dynatrace-bootstrapper",
		Short: "Simple binary for setting up the OneAgent CodeModule in different envs.",
		Long:  "The purpose of the bootstrapper is to copy and configure of the OneAgent CodeModule. If not subcommands are specified, then the 'init' subcommand will run, due to backwards compatibility.",
		RunE:  bootstrapper.RunE,
	}

	bootstrapper.AddFlags(rootCmd)

	rootCmd.AddCommand(bootstrapper.New(), serverless.New())

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
