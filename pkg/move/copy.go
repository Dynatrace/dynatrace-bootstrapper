package move

import (
	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

type copyFunc func(log logr.Logger, fs afero.Afero, from, to string) error

var _ copyFunc = simpleCopy

func simpleCopy(log logr.Logger, fs afero.Afero, from, to string) error {
	log.Info("starting to copy (simple)", "from", from, "to", to)

	err := fsutils.CopyFolder(log, fs, from, to)
	if err != nil {
		log.Error(err, "error moving folder")

		return err
	}

	log.Info("successfully copied (simple)", "from", from, "to", to)

	return nil
}
