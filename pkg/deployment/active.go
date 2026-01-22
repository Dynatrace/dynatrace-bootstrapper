package deployment

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/log"
	"github.com/go-logr/logr"
)

const (
	ActiveLinkName = "active"

	dirPerm755 fs.FileMode = 0o755
)

// CreateActiveSymlinkAtomically creates the `active` symlink pointing to the specified OneAgent target path.
// If the symlink exists, it is updated atomically using a rename operation:
// First, a temporary symlink is created in the work folder and then atomically renamed to the target 'active' symlink.
func CreateActiveSymlinkAtomically(logger logr.Logger, workBaseFolder, agentTargetPath string) error {
	if err := os.MkdirAll(workBaseFolder, dirPerm755); err != nil {
		return fmt.Errorf("failed to create the work base folder: %w", err)
	}

	workFolder, err := os.MkdirTemp(workBaseFolder, "link-work-*")
	if err != nil {
		return fmt.Errorf("failed to create the temporary link work folder: %w", err)
	}

	defer func() {
		if cleanupErr := os.RemoveAll(workFolder); cleanupErr != nil {
			logger.Error(cleanupErr, "failed to cleanup the active link work folder")
		}
	}()

	tmpActiveSymlink := filepath.Join(workFolder, ActiveLinkName)
	agentFolder := filepath.Base(agentTargetPath)

	log.Debug(logger, "Creating a temporary symlink pointing to the target OneAgent folder", "temporary symlink", tmpActiveSymlink, "agent folder", agentFolder)

	if err := os.Symlink(agentFolder, tmpActiveSymlink); err != nil {
		return fmt.Errorf("failed to create the temporary symlink: %w", err)
	}

	activeSymlink := getPathToActiveLink(agentTargetPath)

	log.Debug(logger, "Renaming the temporary symlink to the `active` symlink", "temporary symlink", tmpActiveSymlink,
		"`active` symlink", activeSymlink, "target OneAgent folder", agentFolder)

	if err := os.Rename(tmpActiveSymlink, activeSymlink); err != nil {
		return fmt.Errorf("failed to rename the temporary symlink: %w", err)
	}

	return nil
}

func getPathToActiveLink(agentTargetPath string) string {
	agentDir := filepath.Dir(agentTargetPath)

	return filepath.Join(agentDir, ActiveLinkName)
}
