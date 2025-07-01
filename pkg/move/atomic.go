package move

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

// Atomic provides atomic copy operations by copying to a work directory first, then renaming to target.
func Atomic(work string, copyFunc CopyFunc) CopyFunc {
	return func(log logr.Logger, fs afero.Afero, from, to string) (err error) {
		log.Info("setting up atomic operation", "from", from, "to", to, "work", work)

		err = fs.RemoveAll(work)
		if err != nil {
			log.Error(err, "failed initial cleanup of workdir")

			return err
		}

		err = fs.MkdirAll(work, os.ModePerm)
		if err != nil {
			log.Error(err, "failed to create the base workdir")

			return err
		}

		defer func() {
			if err != nil {
				if cleanupErr := fs.RemoveAll(work); cleanupErr != nil {
					log.Error(err, "failed cleanup of workdir after failure")
				}
			}
		}()

		err = copyFunc(log, fs, from, work)
		if err != nil {
			log.Error(err, "error copying folder")

			return err
		}

		err = fs.Rename(work, to)
		if err != nil {
			log.Error(err, "error moving folder")

			return err
		}

		log.Info("successfully finalized atomic operation", "from", from, "to", to, "work", work)

		return nil
	}
}
