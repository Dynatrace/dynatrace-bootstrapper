package pgc

import (
	"os"
	"path/filepath"

	fs "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
)

const (
	InputFileName              = "declarative.cbor"
	DestinationDeclarativePath = "oneagent/agent/config/declarative.cbor"
)

func GetDestinationFilePath(containerConfigDir string) string {
	return filepath.Join(containerConfigDir, DestinationDeclarativePath)
}

func Configure(log logr.Logger, inputDir, targetDir, containerConfigDir string) error {
	inputFilePath := filepath.Join(inputDir, InputFileName)

	if _, err := os.Stat(inputFilePath); os.IsNotExist(err) {
		return nil
	}

	dstPath := GetDestinationFilePath(containerConfigDir)

	if err := os.MkdirAll(filepath.Dir(dstPath), 0o755); err != nil {
		return err
	}

	log.Info("copying declarative.cbor", "src", inputFilePath, "dst", dstPath)
	return fs.CopyFile(inputFilePath, dstPath)
}
