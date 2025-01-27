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
	sourceFolderFlag = "source"
	targetFolderFlag = "target"
	workFolderFlag   = "work"

	copyTmpFolder = "copy-tmp"
)

var (
	sourceFolder string
	targetFolder string
	workFolder   string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(&sourceFolder, sourceFolderFlag, "", "Base path where to copy the codemodule FROM.")
	_ = cmd.MarkPersistentFlagRequired(sourceFolderFlag)

	cmd.PersistentFlags().StringVar(&targetFolder, targetFolderFlag, "", "Base path where to copy the codemodule TO.")
	_ = cmd.MarkPersistentFlagRequired(targetFolderFlag)

	cmd.PersistentFlags().StringVar(&workFolder, workFolderFlag, "", "(Optional) Base path for a tmp folder, this is where the command will do its work, to make sure the operations are atomic. It must be on the same disk as the target folder.")
}

// Execute moves the contents of a folder to another via copying.
// This could be a simple os.Rename, however that will not work if the source and target are on different disk.
func Execute(fs afero.Afero) error {
	if workFolder != "" {
		return atomicCopy(fs)
	}

	return simpleCopy(fs)
}

func atomicCopy(fs afero.Afero) error {
	logrus.Infof("Starting to copy (atomic) from %s to %s", sourceFolder, targetFolder)

	err := fs.RemoveAll(workFolder)
	if err != nil {
		logrus.Errorf("Failed initial cleanup of workdir: %v", err)

		return err
	}

	err = fs.MkdirAll(workFolder, os.ModePerm)
	if err != nil {
		logrus.Errorf("Failed to create the base workdir: %v", err)

		return err
	}

	defer func() {
		err := fs.RemoveAll(workFolder)
		if err != nil {
			logrus.Errorf("Failed to do cleanup after run: %v", err)
		}
	}()

	tmpFolder := filepath.Join(workFolder, copyTmpFolder)

	err = copyFolder(fs, sourceFolder, tmpFolder)
	if err != nil {
		logrus.Errorf("Error moving folder: %v", err)

		return err
	}

	err = fs.Rename(tmpFolder, targetFolder)
	if err != nil {
		logrus.Errorf("Error finalizing move: %v", err)

		return err
	}

	logrus.Infof("Successfully copied from %s to %s", sourceFolder, targetFolder)

	return nil
}

func simpleCopy(fs afero.Afero) error {
	logrus.Infof("Starting to copy (simple) from %s to %s", sourceFolder, targetFolder)

	err := copyFolder(fs, sourceFolder, targetFolder)
	if err != nil {
		logrus.Errorf("Error moving folder: %v", err)

		return err
	}

	logrus.Infof("Successfully copied from %s to %s", sourceFolder, targetFolder)

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
