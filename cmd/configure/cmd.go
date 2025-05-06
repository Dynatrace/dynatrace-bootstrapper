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
	IsFullstackFlag  = "fullstack"
	TenantFlag       = "tenant"
)

var (
	inputFolder  string
	configFolder string
	installPath  = "/opt/dynatrace/oneagent"
	isFullstack  bool
	tenant       string

	podAttributes       []string
	containerAttributes []string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&inputFolder, InputFolderFlag, "", "(Optional) Base path where to look for the configuration files.")

	cmd.PersistentFlags().StringVar(&configFolder, ConfigFolderFlag, "", "(Optional) Base path where to put the configuration files.")

	cmd.PersistentFlags().StringVar(&installPath, InstallPathFlag, "/opt/dynatrace/oneagent", "(Optional) Base path where the agent binary will be put.")

	cmd.PersistentFlags().StringArrayVar(&containerAttributes, container.Flag, []string{}, "(Optional) Container-specific attributes in JSON format.")

	cmd.PersistentFlags().StringArrayVar(&podAttributes, pod.Flag, []string{}, "(Optional) Pod-specific attributes in key=value format.")

	cmd.PersistentFlags().BoolVar(&isFullstack, IsFullstackFlag, false, "(Optional) Configure the CodeModule to be fullstack.")

	cmd.PersistentFlags().Lookup(IsFullstackFlag).NoOptDefVal = "true"

	cmd.PersistentFlags().StringVar(&tenant, TenantFlag, "", "The name of the tenant that the CodeModule will communicate with. Mandatory in case of --fullstack.")
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

	podAttr, err := pod.ParseAttributes(podAttributes)
	if err != nil {
		return err
	}

	containerAttrs, err := container.ParseAttributes(containerAttributes)
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

		err = conf.Configure(log, fs, containerConfigDir, containerAttr, podAttr, tenant, isFullstack)
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

	log.Info("finished configuration", "config-directory", configFolder, "input-directory", inputFolder)

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
