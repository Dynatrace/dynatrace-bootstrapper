package main

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/version"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func main() {
	cmd := bootstrapper(afero.NewOsFs())

	err := cmd.Execute()
	if err != nil {
		logrus.Errorf("Error during bootstrapping: %v", err)
		os.Exit(1)
	}
}

func bootstrapper(fs afero.Fs) *cobra.Command {
	cmd := &cobra.Command{
		Use:  "dynatrace-bootstrapper",
		RunE: run(fs),
	}

	move.AddFlags(cmd)
	configure.AddFlags(cmd)

	return cmd
}

func run(fs afero.Fs) func(cmd *cobra.Command, _ []string) error {
	return func(cmd *cobra.Command, _ []string) error {
		version.Print()

		err := cmd.ValidateRequiredFlags()
		if err != nil {
			return err
		}

		aferoFs := afero.Afero{
			Fs: fs,
		}

		err = move.Execute(aferoFs)
		if err != nil {
			return err
		}

		return configure.Execute(aferoFs)
	}
}
