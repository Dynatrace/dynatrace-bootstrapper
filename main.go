package main

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/azureappservice"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/k8sinit"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "dynatrace-bootstrapper",
		Short: "Simple binary for setting up the OneAgent CodeModule in different envs.",
		Long:  "The purpose of the bootstrapper is to copy and configure of the OneAgent CodeModule. If not subcommands are specified, then the 'init' subcommand will run, due to backwards compatibility.",
		RunE:  k8sinit.RunE,
	}

	k8sinit.AddFlags(rootCmd)

	rootCmd.AddCommand(
		k8sinit.New(),
		azureappservice.New(),
	)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
