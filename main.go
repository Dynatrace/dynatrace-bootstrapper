package main

import (
	"os"

	bootstrapper "github.com/Dynatrace/dynatrace-bootstrapper/cmd/init"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/serverless"
	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "dynatrace-bootstrapper",
		Short: "Short",
		Long: "Long",
	}

	rootCmd.AddCommand(bootstrapper.New(), serverless.New())

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
