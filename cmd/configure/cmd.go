package configure

import (
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/pod"
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
	InputFolderFlag  = "input-directory"
	ConfigFolderFlag = "config-directory"
	InstallPathFlag  = "install-path"
)

var (
	inputDir    string
	configDir   string
	installPath = "/opt/dynatrace/oneagent"

	podAttributes       []string
	containerAttributes []string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&inputDir, InputFolderFlag, "", "(Optional) Base path where to look for the configuration files.")

	cmd.PersistentFlags().StringVar(&configDir, ConfigFolderFlag, "", "(Optional) Base path where to put the configuration files.")

	cmd.PersistentFlags().StringVar(&installPath, InstallPathFlag, "/opt/dynatrace/oneagent", "(Optional) Base path where the agent binary will be put.")

	cmd.PersistentFlags().StringArrayVar(&containerAttributes, container.Flag, []string{}, "(Optional) Container-specific attributes in JSON format.")

	cmd.PersistentFlags().StringArrayVar(&podAttributes, pod.Flag, []string{}, "(Optional) Pod-specific attributes in key=value format.")
}

func Execute(log logr.Logger, fs afero.Afero, targetDir string) error {
	if configDir == "" || inputDir == "" {
		return nil
	}

	log.Info("starting configuration", "config-directory", configDir, "input-directory", inputDir)

	err := preload.Configure(log, fs, configDir, installPath)
	if err != nil {
		log.Info("failed to configure the ld.so.preload", "config-directory", configDir)

		return err
	}

	podAttr, err := pod.ParseAttributes(podAttributes)
	if err != nil {
		return err
	}

	containerAttrs, err := container.ParseAttributes(containerAttributes)
	if err != nil {
		return err
	}

	for _, containerAttr := range containerAttrs {
		containerConfigDir := filepath.Join(configDir, containerAttr.ContainerName)

		err = pmc.Configure(log, fs, inputDir, targetDir, containerConfigDir, installPath)
		if err != nil {
			log.Info("failed to configure the ruxitagentproc.conf", "config-directory", containerConfigDir)

			return err
		}

		err = metadata.Configure(log, fs, containerConfigDir, podAttr, containerAttr)
		if err != nil {
			log.Info("failed to configure the enrichment files", "config-directory", containerConfigDir)

			return err
		}

		err = endpoint.Configure(log, fs, inputDir, containerConfigDir)
		if err != nil {
			log.Info("failed to configure the endpoint.properties", "config-directory", configDir)

			return err
		}

		err = conf.Configure(log, fs, containerConfigDir, containerAttr, podAttr)
		if err != nil {
			log.Info("failed to configure the container-conf files", "config-directory", containerConfigDir)

			return err
		}

		err = configureFromInputDir(log, fs, containerConfigDir, inputDir)
		if err != nil {
			log.Info("failed to configure container", "config-directory", containerConfigDir)

			return err
		}
	}

	log.Info("finished configuration", "config-directory", configDir, "input-directory", inputDir)

	return nil
}

func configureFromInputDir(log logr.Logger, fs afero.Afero, containerConfigDir, inputDir string) error {
	log.Info("starting to configure the container", "path", containerConfigDir)

	err := curl.Configure(log, fs, inputDir, containerConfigDir)
	if err != nil {
		log.Info("failed to configure the curl options", "config-directory", containerConfigDir)

		return err
	}

	err = ca.Configure(log, fs, inputDir, containerConfigDir)
	if err != nil {
		log.Info("failed to configure the CAs", "config-directory", containerConfigDir)

		return err
	}

	return nil
}
