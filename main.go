package main

import (
	"fmt"
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/version"
	"github.com/spf13/cobra"
)

func main() {
	cmd := new()

	err := cmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func new() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "dynatrace-bootstrapper",
		RunE: base,
	}

	return cmd
}

func base(_ *cobra.Command, _ []string) error {
	version.Print()
	return nil
}
