package pgc

import (
	"os"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
)

const (
	InputFileName              = "declarative.cbor"
	DestinationDeclarativePath = "oneagent/agent/config/" + InputFileName
	FileMode = os.FileMode(0600)
)

func GetDestinationFilePath(containerConfigDir string) string {
	return filepath.Join(containerConfigDir, DestinationDeclarativePath)
}

func Configure(log logr.Logger, inputDir, _ string, containerConfigDir string) error {
	inputFilePath := filepath.Join(inputDir, InputFileName)

	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return nil
	}

	dstPath := GetDestinationFilePath(containerConfigDir)

	if err := os.MkdirAll(filepath.Dir(dstPath), FileMode); err != nil {
		return err
	}

	log.Info("copying declarative.cbor", "src", inputFilePath, "dst", dstPath)

	return fs.CopyFile(inputFilePath, dstPath)
}
