// Package fs provides filesystem utilities for copying and managing files and directories.
package fs

import (
	"io"
	"os"
	"path/filepath"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

// CopyFolder recursively copies a folder from one location to another using the provided afero filesystem.
func CopyFolder(log logr.Logger, fs afero.Fs, from string, to string) error {
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
		fromPath := filepath.Join(from, entry.Name())
		toPath := filepath.Join(to, entry.Name())

		if entry.IsDir() {
			log.V(1).Info("copying directory", "from", fromPath, "to", toPath)

			err = CopyFolder(log, fs, fromPath, toPath)
			if err != nil {
				return err
			}
		} else {
			log.V(1).Info("copying file", "from", fromPath, "to", toPath)

			err = CopyFile(fs, fromPath, toPath)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies a file from sourcePath to destinationPath using the provided afero filesystem.
func CopyFile(fs afero.Fs, sourcePath string, destinationPath string) error {
	sourceFile, err := fs.Open(sourcePath)
	if err != nil {
		return errors.WithStack(err)
	}

	defer func() { _ = sourceFile.Close() }()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	destinationFile, err := fs.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
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
