package move

import (
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs/symlink"
	"github.com/go-logr/logr"
	"github.com/spf13/afero"
)

const (
	InstallerVersionFilePath = "agent/installer.version"
	currentDir               = "agent/bin/current"
)

// CreateCurrentSymlink finds the version of the CodeModule in the `targetDir` (in the installer.version file) and creates a "current" symlink next to it.
// this is needed for the nginx use-case.
func CreateCurrentSymlink(log logr.Logger, fs afero.Afero, targetDir string) error {
	targetCurrentDir := filepath.Join(targetDir, currentDir)

	exists, err := fs.Exists(targetCurrentDir)
	if exists {
		log.Info("the current version dir already exists, skipping symlinking", "current version dir", targetCurrentDir)

		return nil
	} else if err != nil {
		log.Info("failed to check the state of the current version dir", "current version dir", targetCurrentDir)

		return err
	}

	versionFilePath := filepath.Join(targetDir, InstallerVersionFilePath)
	version, err := fs.ReadFile(versionFilePath)

	if err != nil {
		log.Info("failed to get the version from the filesystem", "version-file", versionFilePath)

		return err
	}

	return symlink.Create(log, fs.Fs, string(version), targetCurrentDir)
}
