package configure

import (
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/config_files"
	"github.com/sirupsen/logrus"
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

func Execute(fs afero.Afero) error {
	if configDirectory == "" {
		return nil
	}

	podAttr, err := pod.ParseAttributes()
	if err != nil {
		return err
	}

	containerAttrs, err := container.ParseAttributes()
	if err != nil {
		return err
	}

	logrus.Info("Starting to configure enrichment files")

	for _, containerAttr := range containerAttrs {
		err = config_files.ConfigureEnrichmentFiles(fs, configDirectory, *podAttr, containerAttr.ContainerName)
		if err != nil {
			logrus.Infof("Failed to configure the enrichment files, config-directory: %s", configDirectory)
			return err
		}

		err = config_files.ConfigureContainerConfFile(fs, configDirectory, containerAttr)
		if err != nil {
			logrus.Infof("Failed to configure the container-conf files, config-directory: %s", configDirectory)
			return err
		}

	}

	return nil
}
