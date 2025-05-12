package curl

import (
	"fmt"
	"os"
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	optionsFormatString = `initialConnectRetryMs %s
`
	ConfigPath    = "oneagent/agent/customkeys/curl_options.conf"
	InputFileName = "initial-connect-retry"
)

func Configure(log logr.Logger, fs afero.Afero, inputDir, configDir string) error {
	content, err := getFromFs(fs, inputDir)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("input file not present, skipping curl options configuration", "path", filepath.Join(inputDir, InputFileName))

			return nil
		}

		return err
	}

	configFile := filepath.Join(configDir, ConfigPath)

	log.Info("configuring curl_options.conf", "config-path", configFile)

	return fsutils.CreateFile(fs, configFile, content)
}

func getFromFs(fs afero.Afero, inputDir string) (string, error) {
	inputFile := filepath.Join(inputDir, InputFileName)

	content, err := fs.ReadFile(inputFile)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(optionsFormatString, string(content)), err
}
