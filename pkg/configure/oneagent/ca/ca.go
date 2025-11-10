package ca

import (
	"os"
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
)

const (
	ConfigBasePath     = "oneagent/agent/customkeys"
	ProxyCertsFileName = "custom_proxy.pem"
	CertsFileName      = "custom.pem"

	TrustedCertsInputFile = "trusted.pem"
	AgCertsInputFile      = "activegate.pem"
)

func Configure(log logr.Logger, inputDir, configDir string) error {
	trustedCerts, err := GetFromFs(inputDir, TrustedCertsInputFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	agCerts, err := GetFromFs(inputDir, AgCertsInputFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if agCerts != "" || trustedCerts != "" {
		certFilePath := filepath.Join(configDir, ConfigBasePath, CertsFileName)
		log.Info("creating cert file", "path", certFilePath)

		err := fsutils.CreateFile(certFilePath, agCerts+"\n"+trustedCerts)
		if err != nil {
			return err
		}
	}

	if trustedCerts != "" {
		proxyCertFilePath := filepath.Join(configDir, ConfigBasePath, ProxyCertsFileName)
		log.Info("creating cert file", "path", proxyCertFilePath)

		err := fsutils.CreateFile(proxyCertFilePath, trustedCerts)
		if err != nil {
			return err
		}
	}

	return nil
}

func GetFromFs(inputDir, certFileName string) (string, error) {
	inputFile := filepath.Join(inputDir, certFileName)

	content, err := os.ReadFile(inputFile)
	if err != nil {
		return "", err
	}

	return string(content), err
}
