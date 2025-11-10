package move

import (
	"os"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs/symlink"
	"github.com/go-logr/logr"
)

const (
	InstallerVersionFilePath = "agent/installer.version"
	CurrentDir               = "agent/bin/current"
)

// CreateCurrentSymlink finds the version of the CodeModule in the `targetDir` (in the installer.version file) and creates a "current" symlink in the agent/bin folder that points to the agent/bin/<version> subfolder.
// this is needed for the nginx use-case.
func CreateCurrentSymlink(log logr.Logger, targetDir string) error {
	targetCurrentDir := filepath.Join(targetDir, CurrentDir)

	stat, err := os.Stat(targetCurrentDir)
	if stat != nil {
		log.Info("the current version dir already exists, skipping symlinking", "current version dir", targetCurrentDir)

		return nil
	} else if err != nil && !os.IsNotExist(err) {
		log.Info("failed to check the state of the current version dir", "current version dir", targetCurrentDir)

		return err
	}

	versionFilePath := filepath.Join(targetDir, InstallerVersionFilePath)

	version, err := os.ReadFile(versionFilePath)
	if err != nil {
		log.Info("failed to get the version from the filesystem", "version-file", versionFilePath)

		return err
	}

	return symlink.Create(log, string(version), targetCurrentDir)
}
