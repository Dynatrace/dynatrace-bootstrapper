package conf

import (
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	configPath = "/oneagent/agent/config/container.conf"
)

func Configure(log logr.Logger, fs afero.Afero, configDirectory string, containerAttr container.Attributes, podAttr pod.Attributes) error {
	confContent := fromAttributes(containerAttr, podAttr)

	stringContent, err := confContent.toString()
	if err != nil {
		log.Error(err, "failed to create container conf content", "struct", confContent)

		return err
	}

	configFilePath := filepath.Join(configDirectory, configPath)

	err = fsutils.CreateFile(fs, configFilePath, stringContent)
	if err != nil {
		log.Error(err, "failed to create container conf file", "struct", configFilePath)

		return err
	}

	return nil
}
