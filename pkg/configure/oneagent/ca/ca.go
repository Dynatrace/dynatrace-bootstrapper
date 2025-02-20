package ca

import (
	"os"
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	configBasePath     = "oneagent/agent/customkeys"
	proxyCertsFileName = "custom_proxy.pem"
	certsFileName      = "custom.pem"

	TrustedCertsInputFile = "trusted.pem"
	AgCertsInputFile      = "activegate.pem"
)

func Configure(log logr.Logger, fs afero.Afero, inputDir, configDir string) error {
	trustedCerts, err := getFromFs(fs, inputDir, TrustedCertsInputFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	agCerts, err := getFromFs(fs, inputDir, AgCertsInputFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if agCerts != "" || trustedCerts != "" {
		certFilePath := filepath.Join(configDir, configBasePath, certsFileName)
		log.Info("creating cert file", "path", certFilePath)

		err := fsutils.CreateFile(fs, certFilePath, agCerts+"\n"+trustedCerts)
		if err != nil {
			return err
		}

	}

	if trustedCerts != "" {
		proxyCertFilePath := filepath.Join(configDir, configBasePath, proxyCertsFileName)
		log.Info("creating cert file", "path", proxyCertFilePath)

		err := fsutils.CreateFile(fs, proxyCertFilePath, trustedCerts)
		if err != nil {
			return err
		}

	}

	return nil
}

func getFromFs(fs afero.Afero, inputDir, certFileName string) (string, error) {
	inputFile := filepath.Join(inputDir, certFileName)

	content, err := fs.ReadFile(inputFile)
	if err != nil {
		return "", err
	}

	return string(content), err
}
