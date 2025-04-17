package symlink

import (
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func Create(log logr.Logger, fs afero.Fs, targetDir, symlinkDir string) error {
	// MemMapFs (used for testing) doesn't comply with the Linker interface
	linker, ok := fs.(afero.Linker)
	if !ok {
		log.Info("symlinking not possible", "targetDir", targetDir, "fs", fs)

		return nil
	}

	// Check if the symlink already exists
	if fileInfo, _ := fs.Stat(symlinkDir); fileInfo != nil {
		log.Info("symlink already exists", "location", symlinkDir)

		return nil
	}

	log.Info("creating symlink", "points-to(relative)", targetDir, "location", symlinkDir)

	if err := linker.SymlinkIfPossible(targetDir, symlinkDir); err != nil {
		log.Info("symlinking failed", "source", targetDir)

		return errors.WithStack(err)
	}

	return nil
}
