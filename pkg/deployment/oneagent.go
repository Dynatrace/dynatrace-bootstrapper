package deployment

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/lock"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/log"
	"github.com/go-logr/logr"
)

const (
	allTechValue       = "all" // if set, all technologies will be copied
	deploymentLockFile = "deployment.lock"
)

// DeployOneAgent deploys OneAgent to the target directory using an exclusive file lock to prevent concurrent
// deployments in multi-instance environments.
// The lock file is created in the work base folder, ensuring only one instance performs the deployment at a time.
//
// Returns:
// - bool: true if the OneAgent deployment was performed, false if the deployment was skipped (e.g., OneAgent is already deployed or another instance holds the lock)
// - error: if deployment fails or an error occurs during the deployment process
func DeployOneAgent(logger logr.Logger, sourceBaseFolder, targetBaseFolder, workBaseFolder string, technology string) (bool, error) {
	if err := os.MkdirAll(workBaseFolder, dirPerm755); err != nil {
		return false, fmt.Errorf("error creating work base folder: %w", err)
	}

	lockFilePath := getPathToDeploymentLockFile(workBaseFolder)
	fileLock := lock.New(logger, lockFilePath)

	log.Debug(logger, "Try to acquire the deployment lock file", "path", lockFilePath)

	acquired, err := fileLock.TryAcquire()
	if err != nil {
		return false, fmt.Errorf("failed to acquire the deployment lock: %w", err)
	}

	if !acquired {
		logger.Info("Another instance holds the deployment lock, skipping deployment")

		return false, nil
	}

	defer func() {
		if releaseErr := fileLock.Release(); releaseErr != nil {
			logger.Error(releaseErr, "failed to release the deployment lock file", "path", lockFilePath)
		}
	}()

	// Before deploying, check the status again in case another Bootstrapper instance
	// has finished deployment and removed the lock file since the last check.
	result := CheckAgentDeploymentStatus(sourceBaseFolder, targetBaseFolder)
	if result.Error != nil {
		return false, fmt.Errorf("failed to check OneAgent deployment status, skip deployment: %w", result.Error)
	}

	if result.Status == Deployed {
		log.Debug(logger, "OneAgent is already deployed", "OneAgent version", result.AgentVersion)

		return false, nil
	}

	logger.Info("The deployment lock file has been acquired. Proceeding with OneAgent deployment", "OneAgent version", result.AgentVersion)

	agentFolder := GetAgentFolder(targetBaseFolder, result.AgentVersion)
	if result.Status == NotDeployed {
		// the versioned agent folder does not exist, copy the agent
		err = copyAgent(logger, sourceBaseFolder, agentFolder, workBaseFolder, technology)
		if err != nil {
			return false, fmt.Errorf("failed to deploy OneAgent in the target directory: %w", err)
		}
	}

	// create or update the `active` symlink to point to the newly deployed versioned agent folder
	err = CreateActiveSymlinkAtomically(logger, workBaseFolder, agentFolder)
	if err != nil {
		return false, fmt.Errorf("failed to create `active` symlink in the target directory: %w", err)
	}

	logger.Info("OneAgent has been successfully deployed", "OneAgent version", result.AgentVersion)

	return true, nil
}

// copyAgent atomically copies OneAgent from the source to the destination.
// Creates a temporary folder, copies code modules from the source to the temporary folder,
// sets up the current symlink and then atomically moves the temporary folder to the versioned OneAgent folder.
// Temporary and versioned OneAgent folders must be on the same disk for the atomic move (i.e. renaming).
func copyAgent(log logr.Logger, sourceBaseFolder, versionedAgentFolder, workBaseFolder string, technology string) error {
	if err := os.MkdirAll(workBaseFolder, dirPerm755); err != nil {
		return fmt.Errorf("failed to create the work base folder: %w", err)
	}

	if err := os.MkdirAll(filepath.Dir(versionedAgentFolder), dirPerm755); err != nil {
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

func getPathToDeploymentLockFile(workBaseFolder string) string {
	return filepath.Join(workBaseFolder, deploymentLockFile)
}
