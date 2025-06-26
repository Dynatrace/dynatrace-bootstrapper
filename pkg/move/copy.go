// Package move provides utilities for moving and copying files and folders.
package move

import (
	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
	"golang.org/x/sys/unix"
)

const (
	// defaultUmask is the umask value to use for file operations.
	defaultUmask = 0000
)

// CopyFunc defines a function signature for copying files or directories.
type CopyFunc func(log logr.Logger, fs afero.Afero, from, to string) error

var _ CopyFunc = SimpleCopy

// SimpleCopy copies a folder from one location to another using the provided logger and filesystem.
func SimpleCopy(log logr.Logger, fs afero.Afero, from, to string) error {
	log.Info("starting to copy (simple)", "from", from, "to", to)

	oldUmask := unix.Umask(defaultUmask)
	defer unix.Umask(oldUmask)

	err := fsutils.CopyFolder(log, fs, from, to)
	if err != nil {
		log.Error(err, "error moving folder")

		return err
	}

	log.Info("successfully copied (simple)", "from", from, "to", to)

	return nil
}
