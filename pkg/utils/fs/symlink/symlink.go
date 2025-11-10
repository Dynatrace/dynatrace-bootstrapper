package symlink

import (
	"os"

	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

func Create(log logr.Logger, targetDir, symlinkDir string) error {
	// Check if the symlink already exists
	if fileInfo, _ := os.Stat(symlinkDir); fileInfo != nil {
		log.Info("symlink already exists", "location", symlinkDir)

		return nil
	}

	log.Info("creating symlink", "points-to(relative)", targetDir, "location", symlinkDir)

	if err := os.Symlink(targetDir, symlinkDir); err != nil {
		log.Info("symlinking failed", "source", targetDir)

		return errors.WithStack(err)
	}

	return nil
}
