package conf

import (
	"errors"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/pod"
	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	ConfigPath = "/oneagent/agent/config/container.conf"
)

func Configure(log logr.Logger, fs afero.Afero, configDirectory string, containerAttr container.Attributes, podAttr pod.Attributes, tenant string, isFullstack bool) error {
	log.Info("configuring container.conf", "config-directory", configDirectory)

	if isFullstack {
		log.Info("fullstack flag detected, configuring accordingly", "tenant", tenant)

		if tenant == "" {
			return errors.New("fullstack mode is set, but no tenant was provided")
		}
	}

	confContent := fromAttributes(containerAttr, podAttr, tenant, isFullstack)

	stringContent, err := confContent.toString()
	if err != nil {
		log.Error(err, "failed to create container conf content", "struct", confContent)

		return err
	}

	configFilePath := filepath.Join(configDirectory, ConfigPath)

	err = fsutils.CreateFile(fs, configFilePath, stringContent)
	if err != nil {
		log.Error(err, "failed to create container conf file", "struct", configFilePath)

		return err
	}

	return nil
}
