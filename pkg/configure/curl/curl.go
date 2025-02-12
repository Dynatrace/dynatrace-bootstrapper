package curl

import (
	"fmt"
	"os"
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/spf13/afero"
)

const (
	optionsFormatString = `initialConnectRetryMs %s
`
	configPath = "oneagent/agent/customkeys/curl_options.conf"
	inputFileName  = "initialConnectRetry"

)

func Configure(fs afero.Afero, inputDir string, configDir string) error {
	content, err := getFromFs(fs, inputDir)
	if os.IsNotExist(err) {
		return nil
	} else if err != nil {
		return err
	}


	return createFile(fs, configDir, content)
}

func getFromFs(fs afero.Afero, inputDir string) (string, error) {
	inputFile := filepath.Join(inputDir, inputFileName)

	content, err := fs.ReadFile(inputFile)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf(optionsFormatString, string(content)), err
}

func createFile(fs afero.Afero, configDir, content string) error {
	configFile := filepath.Join(configDir, configPath)

	return fsutils.CreateFile(fs, configFile, content)
}
