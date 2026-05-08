package pgc

import (
	"os"
	"path/filepath"

	fs "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
)

const (
	InputFileName   = "processgroup.json"
	DestinationPath = "agent/conf/processgroup.json"
)

func Configure(log logr.Logger, inputDir, targetDir string) error {
	inputFilePath := filepath.Join(inputDir, InputFileName)

	if _, err := os.Stat(inputFilePath); err != nil {
		if os.IsNotExist(err) {
			log.Info("Input file not present, skipping processgroup.json configuration", "path", inputFilePath)

			return nil
		}

		return err
	}

	dstPath := filepath.Join(targetDir, DestinationPath)

	if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return err
	}

	log.Info("writing processgroup.json", "destination", dstPath)

	return fs.CopyFile(inputFilePath, dstPath)
}
