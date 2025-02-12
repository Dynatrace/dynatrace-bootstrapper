package move

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func atomic(copy copyFunc) copyFunc {
	return func(fs afero.Afero, from, to string) error {
		logrus.Infof("Setting up atomic operation from %s to %s", from, to)

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

		err = copy(fs, from, workFolder)
		if err != nil {
			logrus.Errorf("Error copying folder: %v", err)

			return err
		}

		err = fs.Rename(workFolder, to)
		if err != nil {
			logrus.Errorf("Error moving folder: %v", err)

			return err
		}

		logrus.Infof("Successfully finalized atomic operation from %s to %s", workFolder, to)

		return nil
	}
}
