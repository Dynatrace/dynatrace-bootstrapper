package configure

import (
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/enrichment/endpoint"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/enrichment/metadata"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/ca"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/conf"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/curl"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/preload"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	inputFolderFlag  = "input-directory"
	configFolderFlag = "config-directory"
	installPathFlap  = "install-path"
)

var (
	inputFolder  string
	configFolder string
	installPath  = "/opt/dynatrace/oneagent"
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&inputFolder, inputFolderFlag, "", "(Optional) Base path where to look for the configuration files.")

	cmd.PersistentFlags().StringVar(&configFolder, configFolderFlag, "", "(Optional) Base path where to put the configuration files.")

	cmd.PersistentFlags().StringVar(&installPath, installPathFlap, "/opt/dynatrace/oneagent", "(Optional) Base path where the agent binary will be put.")

	container.AddFlags(cmd)
	pod.AddFlags(cmd)
}

func Execute(log logr.Logger, fs afero.Afero, targetDir string) error {
	if configFolder == "" || inputFolder == "" {
		return nil
	}

	log.Info("starting configuration", "config-directory", configFolder, "input-directory", inputFolder)

	err := preload.Configure(log, fs, configFolder, installPath)
	if err != nil {
		log.Info("failed to configure the ld.so.preload", "config-directory", configFolder)

		return err
	}

	err = pmc.Configure(log, fs, inputFolder, targetDir)
	if err != nil {
		log.Info("failed to configure the ruxitagentproc.conf", "config-directory", configFolder)

		return err
	}

	podAttr, err := pod.ParseAttributes()
	if err != nil {
		return err
	}

	containerAttrs, err := container.ParseAttributes()
	if err != nil {
		return err
	}

	for _, containerAttr := range containerAttrs {
		containerConfigDir := filepath.Join(configFolder, containerAttr.ContainerName)

		err = metadata.Configure(log, fs, containerConfigDir, podAttr, containerAttr)
		if err != nil {
			log.Info("failed to configure the enrichment files", "config-directory", containerConfigDir)

			return err
		}

		err = endpoint.Configure(log, fs, inputFolder, containerConfigDir)
		if err != nil {
			log.Info("failed to configure the endpoint.properties", "config-directory", configFolder)

			return err
		}

		err = conf.Configure(log, fs, containerConfigDir, containerAttr, podAttr)
		if err != nil {
			log.Info("failed to configure the container-conf files", "config-directory", containerConfigDir)

			return err
		}

		err = configureFromInputDir(log, fs, containerConfigDir, inputFolder)
		if err != nil {
			log.Info("failed to configure container", "config-directory", containerConfigDir)

			return err
		}
	}

	log.Info("finished to configuration", "config-directory", configFolder, "input-directory", inputFolder)

	return nil
}

func configureFromInputDir(log logr.Logger, fs afero.Afero, configDir, inputDir string) error {
	log.Info("starting to configure the container", "path", configFolder)

	err := curl.Configure(log, fs, inputDir, configDir)
	if err != nil {
		log.Info("failed to configure the curl options", "config-directory", configFolder)

		return err
	}

	err = ca.Configure(log, fs, inputDir, configDir)
	if err != nil {
		log.Info("failed to configure the CAs", "config-directory", configFolder)

		return err
	}

	return nil
}
