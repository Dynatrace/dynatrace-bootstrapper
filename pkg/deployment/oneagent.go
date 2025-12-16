package deployment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/go-logr/logr"
)

const (
	allTechValue = "all" // if set, all technologies will be copied
)

// CopyAgent atomically copies OneAgent from the source to the destination.
// Creates a temporary folder, copies OneAgent from the source to the temporary folder,
// sets up the current symlink and then atomically moves the temporary folder to the versioned OneAgent folder.
// Temporary and versioned OneAgent folders must be on the same disk for the atomic move (i.e. renaming).
func CopyAgent(log logr.Logger, sourceBaseFolder, versionedAgentFolder, workBaseFolder string, technology string) error {
	if err := os.MkdirAll(workBaseFolder, 0755); err != nil {
		return fmt.Errorf("failed to create the work base folder: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(versionedAgentFolder), 0755); err != nil {
		return fmt.Errorf("failed to create the target folder: %w", err)
	}

	copyFunc := move.SimpleCopy
	if technology != "" && strings.TrimSpace(technology) != allTechValue {
		copyFunc = move.CopyByTechnologyWrapper(technology)
	}

	workFolder, err := os.MkdirTemp(workBaseFolder, "copy-work-*")
	if err != nil {
		return fmt.Errorf("failed to create the temporary copy work folder: %w", err)
	}
	defer func() {
		if cleanupErr := os.RemoveAll(workFolder); cleanupErr != nil {
			log.Error(cleanupErr, "failed to cleanup the copy work folder")
		}
	}()

	copyFunc = move.CreateCurrentSymlinkOnCopy(copyFunc)
	copyFunc = move.Atomic(workFolder, copyFunc)

	return copyFunc(log, sourceBaseFolder, versionedAgentFolder)
}
