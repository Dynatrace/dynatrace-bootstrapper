package metadata

import (
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	JsonFilePath       = "enrichment/dt_metadata.json"
	PropertiesFilePath = "enrichment/dt_metadata.properties"
)

func Configure(log logr.Logger, fs afero.Afero, configDirectory string, podAttr pod.Attributes, containerAttr container.Attributes) error {
	confContent := fromAttributes(containerAttr, podAttr)

	log.V(1).Info("format content into a raw form", "struct", confContent)

	confJson, err := confContent.toJson()
	if err != nil {
		return err
	}

	jsonFilePath := filepath.Join(configDirectory, JsonFilePath)

	err = fsutils.CreateFile(fs, jsonFilePath, string(confJson))
	if err != nil {
		log.Error(err, "failed to create metadata-enrichment properties file", "struct", jsonFilePath)

		return err
	}

	confProperties, err := confContent.toProperties()
	if err != nil {
		return err
	}

	propsFilePath := filepath.Join(configDirectory, PropertiesFilePath)

	err = fsutils.CreateFile(fs, propsFilePath, confProperties)
	if err != nil {
		log.Error(err, "failed to create metadata-enrichment properties file", "struct", propsFilePath)

		return err
	}

	return nil
}
