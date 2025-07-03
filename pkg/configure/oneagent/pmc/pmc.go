package pmc

import (
	"os"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	InputFileName = "ruxitagentproc.json"

	SourceRuxitAgentProcPath      = "agent/conf/ruxitagentproc.conf"
	DestinationRuxitAgentProcPath = "oneagent/agent/config/ruxitagentproc.conf"
)

// GetSourceRuxitAgentProcFilePath returns the full path to the source ruxit agent process configuration file.
func GetSourceRuxitAgentProcFilePath(targetDir string) string {
	return filepath.Join(targetDir, SourceRuxitAgentProcPath)
}

// GetDestinationRuxitAgentProcFilePath returns the full path to the destination ruxit agent process configuration file.
func GetDestinationRuxitAgentProcFilePath(configDir string) string {
	return filepath.Join(configDir, DestinationRuxitAgentProcPath)
}

// Configure applies the configuration from the input file to the ruxit agent process configuration.
func Configure(log logr.Logger, fs afero.Afero, inputDir, targetDir, configDir, installPath string) error {
	inputFilePath := filepath.Join(inputDir, InputFileName)

	inputFile, err := fs.Open(inputFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Info("Input file not present, skipping ruxitagentproc.conf configuration", "path", inputFilePath)

			return nil
		}

		log.Info("failed to input file", "path", inputFilePath)

		return err
	}

	defer func() { _ = inputFile.Close() }()

	conf, err := ruxit.FromJSON(inputFile)
	if err != nil {
		log.Info("failed to unmarshal the input file", "path", inputFilePath)

		return err
	}

	conf.InstallPath = &installPath

	srcPath := GetSourceRuxitAgentProcFilePath(targetDir)
	dstPath := GetDestinationRuxitAgentProcFilePath(configDir)

	log.Info("creating ruxitagentproc.conf", "source", srcPath, "destination", dstPath)

	return Create(log, fs, srcPath, dstPath, conf)
}
