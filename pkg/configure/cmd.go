package configure

import (
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	configDirFlag = "config-directory"
)

var (
	configDirectory string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&configDirectory, configDirFlag, "", "(Optional) Path where enrichment/configuration files will be written.")

	container.AddFlags(cmd)
	pod.AddFlags(cmd)

}

func Execute(fs afero.Afero, to string) error {
	if configDirectory == "" {
		return nil
	}

	podAttr, err := pod.ParseAttributes()
	if err != nil {
		return err
	}

	containerAttr, err := container.ParseAttributes()
	if err != nil {
		return err
	}

	return nil
}
