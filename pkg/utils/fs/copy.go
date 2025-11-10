package fs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

func CopyFolder(log logr.Logger, from string, to string) error {
	fromInfo, err := os.Stat(from)
	if err != nil {
		return errors.WithStack(err)
	}

	if !fromInfo.IsDir() {
		return errors.Errorf("%s is not a directory", from)
	}

	err = os.MkdirAll(to, fromInfo.Mode())
	if err != nil {
		return errors.WithStack(err)
	}

	entries, err := os.ReadDir(from)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, entry := range entries {
		fromPath := filepath.Join(from, entry.Name())
		toPath := filepath.Join(to, entry.Name())

		if entry.IsDir() {
			log.V(1).Info("copying directory", "from", fromPath, "to", toPath)

			err = CopyFolder(log, fromPath, toPath)
			if err != nil {
				return err
			}
		} else {
			log.V(1).Info("copying file", "from", fromPath, "to", toPath)

			err = CopyFile(fromPath, toPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func CopyFile(sourcePath string, destinationPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() { _ = sourceFile.Close() }()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	destinationFile, err := os.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() { _ = destinationFile.Close() }()

	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return errors.WithStack(err)
	}

	err = destinationFile.Sync()
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
