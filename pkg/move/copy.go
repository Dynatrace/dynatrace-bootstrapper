package move

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

type copyFunc func(fs afero.Afero, from, to string) error

var _ copyFunc = simpleCopy

func simpleCopy(fs afero.Afero, from, to string) error {
	logrus.Infof("Starting to copy (simple) from %s to %s", from, to)

	err := copyFolder(fs, from, to)
	if err != nil {
		logrus.Errorf("Error moving folder: %v", err)

		return err
	}

	logrus.Infof("Successfully copied from %s to %s", from, to)

	return nil
}

func copyFolder(fs afero.Fs, from, to string) error {
	fromInfo, err := fs.Stat(from)
	if err != nil {
		return errors.WithStack(err)
	}

	if !fromInfo.IsDir() {
		return errors.Errorf("%s is not a directory", from)
	}

	err = fs.MkdirAll(to, fromInfo.Mode())
	if err != nil {
		return errors.WithStack(err)
	}

	entries, err := afero.ReadDir(fs, from)
	if err != nil {
		return errors.WithStack(err)
	}

	for _, entry := range entries {
		toPath := filepath.Join(from, entry.Name())
		fromPath := filepath.Join(to, entry.Name())

		if entry.IsDir() {
			logrus.Infof("Copying directory %s to %s", toPath, fromPath)

			err = copyFolder(fs, toPath, fromPath)
			if err != nil {
				return err
			}
		} else {
			logrus.Infof("Copying file %s to %s", toPath, fromPath)

			err = copyFile(fs, toPath, fromPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(fs afero.Fs, from string, to string) error {
	sourceFile, err := fs.Open(from)
	if err != nil {
		return errors.WithStack(err)
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	destinationFile, err := fs.OpenFile(to, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		return errors.WithStack(err)
	}

	defer destinationFile.Close()

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
