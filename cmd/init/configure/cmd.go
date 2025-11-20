package configure

import (
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/init/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/init/configure/attributes/pod"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/enrichment/endpoint"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/enrichment/metadata"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/ca"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/conf"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/curl"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/preload"
	"github.com/go-logr/logr"
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
	inputDir    string
	configDir   string
	installPath string
	isFullstack bool
	tenant      string

	podAttributes       []string
	containerAttributes []string
)

func AddFlags(cmd *cobra.Command) {
	// common
	cmd.Flags().StringVar(&inputDir, InputFolderFlag, "", "(Optional) Base path where to look for the configuration files.")
	cmd.Flags().StringVar(&configDir, ConfigFolderFlag, "", "(Optional) Base path where to put the configuration files.")
	cmd.Flags().StringArrayVar(&containerAttributes, container.Flag, []string{}, "(Optional) Container-specific attributes in JSON format.")
	cmd.Flags().StringArrayVar(&podAttributes, pod.Flag, []string{}, "(Optional) Pod-specific attributes in key=value format.")

	// oneagent
	cmd.Flags().StringVar(&installPath, InstallPathFlag, "/opt/dynatrace/oneagent", "(Optional) Base path where the agent binary will be put.")
	cmd.Flags().BoolVar(&isFullstack, IsFullstackFlag, false, "(Optional) Configure the CodeModule to be fullstack.")
	cmd.Flags().StringVar(&tenant, TenantFlag, "", "The name of the tenant that the CodeModule will communicate with. Mandatory in case of --fullstack.")

	cmd.Flags().Lookup(IsFullstackFlag).NoOptDefVal = "true"
}

func SetupOneAgent(log logr.Logger, targetDir string) error {
	if configDir == "" || inputDir == "" {
		return nil
	}

	log.Info("starting configuration", "config-directory", configDir, "input-directory", inputDir)

	err := preload.Configure(log, configDir, installPath)
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
		log.Info("starting to configure the container", "path", containerConfigDir)

		err = pmc.Configure(log, inputDir, targetDir, containerConfigDir, installPath)
		if err != nil {
			log.Info("failed to configure the ruxitagentproc.conf", "config-directory", containerConfigDir)

			return err
		}

		err = conf.Configure(log, containerConfigDir, containerAttr, podAttr, tenant, isFullstack)
		if err != nil {
			log.Info("failed to configure the container-conf files", "config-directory", containerConfigDir)

			return err
		}

		err = configureFromInputDir(log, containerConfigDir, inputDir)
		if err != nil {
			log.Info("failed to configure container", "config-directory", containerConfigDir)

			return err
		}
	}

	log.Info("finished oneagent configuration", "config-directory", configDir, "input-directory", inputDir)

	return nil
}

func configureFromInputDir(log logr.Logger, containerConfigDir, inputDir string) error {
	err := curl.Configure(log, inputDir, containerConfigDir)
	if err != nil {
		log.Info("failed to configure the curl options", "config-directory", containerConfigDir)

		return err
	}

	err = ca.Configure(log, inputDir, containerConfigDir)
	if err != nil {
		log.Info("failed to configure the CAs", "config-directory", containerConfigDir)

		return err
	}

	return nil
}

func EnrichWithMetadata(log logr.Logger) error {
	if configDir == "" || inputDir == "" {
		return nil
	}

	log.Info("starting enrichment", "config-directory", configDir, "input-directory", inputDir)

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
		log.Info("starting to enrich the container", "path", containerConfigDir)

		err = endpoint.Configure(log, inputDir, containerConfigDir)
		if err != nil {
			log.Info("failed to configure the endpoint.properties", "config-directory", configDir)

			return err
		}

		err = metadata.Configure(log, containerConfigDir, podAttr, containerAttr)
		if err != nil {
			log.Info("failed to configure the enrichment files", "config-directory", containerConfigDir)

			return err
		}
	}

	log.Info("finished enrichment", "config-directory", configDir, "input-directory", inputDir)

	return nil
}
