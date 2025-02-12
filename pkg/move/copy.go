package move

import (

	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
)

type copyFunc func(fs afero.Afero, from, to string) error

var _ copyFunc = simpleCopy

func simpleCopy(fs afero.Afero, from, to string) error {
	logrus.Infof("Starting to copy (simple) from %s to %s", from, to)

	err := fsutils.CopyFolder(fs, from, to)
	if err != nil {
		logrus.Errorf("Error moving folder: %v", err)

		return err
	}

	logrus.Infof("Successfully copied from %s to %s", from, to)

	return nil
}
