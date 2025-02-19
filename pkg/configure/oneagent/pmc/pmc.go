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
)

func Configure(log logr.Logger, fs afero.Afero, inputDir, targetDir string) error {
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

	defer inputFile.Close()

	conf, err := ruxit.FromJson(inputFile)
	if err != nil {
		log.Info("failed to unmarshal the input file", "path", inputFilePath)

		return err
	}

	return UpdateInPlace(log, fs, targetDir, conf)
}
