package move

import (
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

const (
	sourceFlag = "source"
	targetFlag = "target"
)

var (
	source string
	target string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&source, sourceFlag, "", "Base path where to copy the codemodule FROM.")
	_ = cmd.MarkPersistentFlagRequired(sourceFlag)

	cmd.PersistentFlags().StringVar(&target, targetFlag, "", "Base path where to copy the codemodule TO.")
	_ = cmd.MarkPersistentFlagRequired(targetFlag)
}

// Execute moves the contents of a folder to another via copying.
// This could be a simple os.Rename, however that will not work if the source and target are on different disk.
func Execute(fs afero.Afero) error {
	logrus.Infof("Starting to copy from %s to %s", source, target)

	tmpLocation := filepath.Join(filepath.Dir(target), "tmp")

	err := copyFolder(fs, source, tmpLocation)
	if err != nil {
		logrus.Errorf("Error moving folder: %v", err)

		return err
	}

	err = fs.Rename(tmpLocation, target)
	if err != nil {
		logrus.Errorf("Error finalizing move: %v", err)

		return err
	}

	logrus.Infof("Successfully copied from %s to %s", source, target)

	return nil
}

func copyFolder(fs afero.Fs, from string, to string) error {
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

func copyFile(fs afero.Fs, sourcePath string, destinationPath string) error {
	sourceFile, err := fs.Open(sourcePath)
	if err != nil {
		return errors.WithStack(err)
	}
	defer sourceFile.Close()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		return errors.WithStack(err)
	}

	destinationFile, err := fs.OpenFile(destinationPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
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
