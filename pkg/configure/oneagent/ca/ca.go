// Package ca provides configuration for OneAgent certificate authority files.
package ca

import (
	"os"
	"path/filepath"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

// ConfigBasePath is the base path for custom key configuration.
// ProxyCertsFileName is the file name for proxy certificates.
// CertsFileName is the file name for agent certificates.
// TrustedCertsInputFile is the input file name for trusted certificates.
// AgCertsInputFile is the input file name for ActiveGate certificates.
const (
	ConfigBasePath     = "oneagent/agent/customkeys"
	ProxyCertsFileName = "custom_proxy.pem"
	CertsFileName      = "custom.pem"

	TrustedCertsInputFile = "trusted.pem"
	AgCertsInputFile      = "activegate.pem"
)

// Configure sets up the certificate files for OneAgent.
func Configure(log logr.Logger, fs afero.Afero, inputDir, configDir string) error {
	trustedCerts, err := GetFromFs(fs, inputDir, TrustedCertsInputFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	agCerts, err := GetFromFs(fs, inputDir, AgCertsInputFile)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	if agCerts != "" || trustedCerts != "" {
		certFilePath := filepath.Join(configDir, ConfigBasePath, CertsFileName)
		log.Info("creating cert file", "path", certFilePath)

		err := fsutils.CreateFile(fs, certFilePath, agCerts+"\n"+trustedCerts)
		if err != nil {
			return err
		}
	}

	if trustedCerts != "" {
		proxyCertFilePath := filepath.Join(configDir, ConfigBasePath, ProxyCertsFileName)
		log.Info("creating cert file", "path", proxyCertFilePath)

		err := fsutils.CreateFile(fs, proxyCertFilePath, trustedCerts)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetFromFs reads a certificate file from the filesystem.
func GetFromFs(fs afero.Afero, inputDir, certFileName string) (string, error) {
	inputFile := filepath.Join(inputDir, certFileName)

	content, err := fs.ReadFile(inputFile)
	if err != nil {
		return "", err
	}

	return string(content), err
}
