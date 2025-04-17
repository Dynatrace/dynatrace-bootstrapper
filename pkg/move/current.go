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

// CreateCurrentSymlink finds a the version of the CodeModule in the `targetDir` and creates a "current" symlink next to it.
// this is needed for the nginx use-case.
func CreateCurrentSymlink(log logr.Logger, fs afero.Afero, targetDir string) error {
	versionFilePath := filepath.Join(targetDir, InstallerVersionFilePath)
	version, err := fs.ReadFile(versionFilePath)

	if err != nil {
		log.Info("failed to get the version from the filesystem", "version-file", versionFilePath)

		return err
	}

	targetBinDir := filepath.Join(targetDir, currentDir)

	return symlink.Create(log, fs.Fs, string(version), targetBinDir)
}
