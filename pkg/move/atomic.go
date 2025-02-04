package move

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

func atomic(copy copyFunc) copyFunc {
	return func(fs afero.Afero) error {
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

		//tmpFolder := filepath.Join(workFolder, copyTmpFolder)

		err = copy(fs)
		if err != nil {
			logrus.Errorf("Error moving folder: %v", err)

			return err
		}

		//err = fs.Rename(tmpFolder, targetFolder)
		//if err != nil {
		//	logrus.Errorf("Error finalizing move: %v", err)

		//			return err
		//		}

		logrus.Infof("Successfully copied from %s to %s", sourceFolder, targetFolder)

		return nil
	}
}
