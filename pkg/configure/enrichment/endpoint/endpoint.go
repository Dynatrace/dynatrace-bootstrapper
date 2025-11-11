package endpoint

import (
	"os"
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
)

const (
	ConfigBasePath = "enrichment/endpoint"
	InputFileName  = "endpoint.properties"
)

func Configure(log logr.Logger, inputDir, configDir string) error {
	properties, err := getFromFs(inputDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("input file not present, skipping endpoint.properties configuration", "path", filepath.Join(inputDir, InputFileName))

			return nil
		}

		return err
	}

	propertiesFileName := filepath.Join(configDir, ConfigBasePath, InputFileName)

	err = fsutils.CreateFile(propertiesFileName, properties)
	if err != nil {
		return err
	}

	return nil
}

func getFromFs(inputDir string) (string, error) {
	inputFile := filepath.Join(inputDir, InputFileName)

	content, err := os.ReadFile(inputFile)
	if err != nil {
		return "", err
	}

	return string(content), err
}
