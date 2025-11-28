package main

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/k8sinit"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:                "dynatrace-bootstrapper",
		Short:              "Simple binary for setting up the OneAgent CodeModule in different envs.",
		Long:               "The purpose of the bootstrapper is to copy and configure the OneAgent CodeModule. If no subcommand is specified, the 'k8s-init' subcommand will run for backward compatibility.",
		FParseErrWhitelist: cobra.FParseErrWhitelist{UnknownFlags: true},
		RunE:               k8sinit.RunE,
	}

	// For backward compatibility, add the k8s-init subcommand flags as the default if no subcommand is specified
	k8sinit.AddFlags(rootCmd)

	rootCmd.AddCommand(
		k8sinit.New(),
	)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
