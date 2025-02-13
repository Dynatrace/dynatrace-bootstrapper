package configure

import (
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/ca"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/curl"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/pmc"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/preload"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	inputFolderFlag  = "input-directory"
	configFolderFlag = "config-directory"
)

var (
	inputFolder  string
	configFolder string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&inputFolder, inputFolderFlag, "", "(Optional) Base path where to look for the configuration files.")

	cmd.PersistentFlags().StringVar(&configFolder, configFolderFlag, "", "(Optional) Base path where to put the configuration files.")
}

// Execute moves the contents of a folder to another via copying.
// This could be a simple os.Rename, however that will not work if the source and target are on different disk.
func Execute(fs afero.Afero, targetDir string) error {
	if configFolder == "" || inputFolder == "" {
		return nil
	}

	logrus.Infof("Starting to configure the CodeModule, config-directory: %s, input-directory: %s", configFolder, inputFolder)

	err := preload.Configure(fs, configFolder)
	if err != nil {
		logrus.Infof("Failed to configure the ld.so.preload, config-directory: %s", configFolder)

		return err
	}

	err = pmc.Configure(fs, inputFolder, configFolder, targetDir)
	if err != nil {
		logrus.Infof("Failed to configure the ruxitagentproc.conf, config-directory: %s", configFolder)

		return err
	}


	// TODO: Here comes the part to parse what container's we have (--attribute-container, --attribute)

	containers := []string{"testing"}
	for _, container := range containers {

		// TODO: add container.conf
		// TODO: add enrichment files

		err = configureFromInputDir(fs, configFolder, inputFolder, container)
		if err != nil {
			logrus.Infof("Failed to configure container, config-directory: %s, container-name: %s", configFolder, container)

			return err
		}
	}


	return nil
}

func configureFromInputDir(fs afero.Afero, configDir, inputDir, containerName string) error {
	logrus.Infof("Starting to configure the container, config-directory: %s, container-name: %s", configDir, containerName)

	containerConfigDir := filepath.Join(configDir, containerName)

	err := curl.Configure(fs, inputDir, containerConfigDir)
	if err != nil {
		logrus.Infof("Failed to configure the curl options, config-directory: %s", configFolder)

		return err
	}

	err = ca.Configure(fs, inputDir, containerConfigDir)
	if err != nil {
		logrus.Infof("Failed to configure the CAs, config-directory: %s", configFolder)

		return err
	}

	logrus.Infof("Finished to configure the CodeModule, config-directory: %s, input-directory: %s", configFolder, inputDir)

	return nil
}
