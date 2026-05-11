package pgc

import (
	"os"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
)

const (
	InputFileName              = "declarative.cbor"
	DestinationDeclarativePath = "oneagent/agent/config/declarative.cbor"
)

func GetDestinationFilePath(containerConfigDir string) string {
	return filepath.Join(containerConfigDir, DestinationDeclarativePath)
}

func Configure(log logr.Logger, inputDir, containerConfigDir string) error {
	inputFilePath := filepath.Join(inputDir, InputFileName)

	srcInfo, err := os.Stat(inputFilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return err
	}

	dstPath := GetDestinationFilePath(containerConfigDir)

	if err := os.MkdirAll(filepath.Dir(dstPath), os.ModePerm); err != nil {
		return err
	}

	log.Info("copying declarative.cbor", "src", inputFilePath, "dst", dstPath)

	if err := fs.CopyFile(inputFilePath, dstPath); err != nil {
		return err
	}

	return os.Chmod(dstPath, srcInfo.Mode().Perm())
}
