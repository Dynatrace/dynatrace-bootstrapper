package endpoint

import (
	"os"
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	configBasePath = "enrichment/endpoint"
	InputFileName  = "endpoint.properties"
)

func Configure(log logr.Logger, fs afero.Afero, inputDir, configDir string) error {
	properties, err := getFromFs(fs, inputDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("input file not present, skipping endpoint.properties configuration", "path", filepath.Join(inputDir, InputFileName))

			return nil
		}

		return err
	}

	propertiesFileName := filepath.Join(configDir, configBasePath, InputFileName)

	err = fsutils.CreateFile(fs, propertiesFileName, properties)
	if err != nil {
		return err
	}

	return nil
}

func getFromFs(fs afero.Afero, inputDir string) (string, error) {
	inputFile := filepath.Join(inputDir, InputFileName)

	content, err := fs.ReadFile(inputFile)
	if err != nil {
		return "", err
	}

	return string(content), err
}
