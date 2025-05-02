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
	DestinationRuxitAgentProcPath = "oneagent/config/ruxitagentproc.conf"
)

func GetSourceRuxitAgentProcPath(targetDir string) string {
	return filepath.Join(targetDir, SourceRuxitAgentProcPath)
}

func GetDestinationRuxitAgentProcPath(configDir string) string {
	return filepath.Join(configDir, DestinationRuxitAgentProcPath)
}

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

	conf, err := ruxit.FromJson(inputFile)
	if err != nil {
		log.Info("failed to unmarshal the input file", "path", inputFilePath)

		return err
	}

	conf.InstallPath = &installPath

	source := GetSourceRuxitAgentProcPath(targetDir)
	destination := GetDestinationRuxitAgentProcPath(configDir)

	log.Info("creating ruxitagentproc.conf", "source", source, "destination", destination)

	return Create(log, fs, source, destination, conf)
}
