package move

import (
	"os"

	"github.com/go-logr/logr"
)

func Atomic(work string, copyFunc CopyFunc) CopyFunc {
	return func(log logr.Logger, from, to string) (err error) {
		log.Info("setting up atomic operation", "from", from, "to", to, "work", work)

		err = os.RemoveAll(work)
		if err != nil {
			log.Error(err, "failed initial cleanup of workdir")

			return err
		}

		err = os.MkdirAll(work, os.ModePerm)
		if err != nil {
			log.Error(err, "failed to create the base workdir")

			return err
		}

		defer func() {
			if err != nil {
				if cleanupErr := os.RemoveAll(work); cleanupErr != nil {
					log.Error(cleanupErr, "failed cleanup of workdir after failure")
				}
			}
		}()

		err = copyFunc(log, from, work)
		if err != nil {
			log.Error(err, "error copying folder")

			return err
		}

		err = os.Rename(work, to)
		if err != nil {
			log.Error(err, "error moving folder")

			return err
		}

		log.Info("successfully finalized atomic operation", "from", from, "to", to, "work", work)

		return nil
	}
}
