package deployment

import (
	"io/fs"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/lock"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	dirPerm000 fs.FileMode = 0o000
	dirPerm700 fs.FileMode = 0o700
)

func TestCopyAgent(t *testing.T) {
	t.Run("Successful copy from Source to Target", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		workBaseDir := t.TempDir()
		targetBaseDir := t.TempDir()
		agentFolder := GetAgentFolder(targetBaseDir, agentVersion)
		err := copyAgent(logger, sourceBaseDir, agentFolder, workBaseDir, allTechValue)
		require.NoError(t, err)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		// the `active` symlink must be missing because it is created after the copy operation
		require.Equal(t, LinkMissing, result.Status)
		require.Equal(t, agentVersion, result.AgentVersion)

		// the `current` symlink must exist because it is created during the copy operation
		currentSymlinkPath := filepath.Join(agentFolder, move.CurrentDir)
		stat, err := os.Lstat(currentSymlinkPath)
		require.NoError(t, err)
		require.NotEqual(t, 0, stat.Mode()&os.ModeSymlink, "current should be a symlink type")
	})

	t.Run("Cannot copy due to a permission error in the target directory", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		err := os.Chmod(targetBaseDir, dirPerm000)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, os.Chmod(targetBaseDir, dirPerm755)) // restore permissions on exit to allow cleanup of the target directory
		}()

		workBaseDir := t.TempDir()
		agentFolder := GetAgentFolder(targetBaseDir, agentVersion)
		err = copyAgent(logger, sourceBaseDir, agentFolder, workBaseDir, allTechValue)
		require.ErrorIs(t, err, syscall.EACCES)

		expectedLog := `failed to create the target folder: mkdir .+: permission denied`
		require.Regexp(t, expectedLog, err.Error())
	})
}

func TestDeployOneAgent(t *testing.T) {
	t.Run("Successfully deploy OneAgent when status is NotDeployed", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.Equal(t, NotDeployed, result.Status)

		workBaseDir := t.TempDir()
		deployed, err := DeployOneAgent(logger, sourceBaseDir, targetBaseDir, workBaseDir, allTechValue)
		require.NoError(t, err)
		require.True(t, deployed)

		// verify that OneAgent is deployed
		result = CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.NoError(t, result.Error)
		require.Equal(t, Deployed, result.Status)
		require.Equal(t, agentVersion, result.AgentVersion)

		// verify the lock file is removed
		lockFilePath := getPathToDeploymentLockFile(workBaseDir)
		_, err = os.Stat(lockFilePath)
		assert.True(t, os.IsNotExist(err), "lock file should be removed after the deployment")
	})

	t.Run("Successfully creates `active` symlink when status is LinkMissing", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		// setup target directory with agent folder but no active symlink
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, "")

		// verify the current status is LinkMissing
		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.Equal(t, LinkMissing, result.Status)

		workBaseDir := t.TempDir()
		deployed, err := DeployOneAgent(logger, sourceBaseDir, targetBaseDir, workBaseDir, allTechValue)
		require.NoError(t, err)
		require.True(t, deployed)

		// verify that OneAgent is deployed
		result = CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.NoError(t, result.Error)
		require.Equal(t, Deployed, result.Status)
		require.Equal(t, agentVersion, result.AgentVersion)

		// verify the lock file is removed
		lockFilePath := getPathToDeploymentLockFile(workBaseDir)
		_, err = os.Stat(lockFilePath)
		assert.True(t, os.IsNotExist(err), "lock file should be removed after the deployment")
	})

	t.Run("Skips deployment when lock already held by another instance", func(t *testing.T) {
		logger, logsObserver := tests.NewTestLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		workBaseDir := t.TempDir()
		lockFilePath := getPathToDeploymentLockFile(workBaseDir)

		// acquire the lock file to simulate another instance holding the lock
		fileLock := lock.New(logger, lockFilePath)
		acquired, err := fileLock.TryAcquire()
		require.NoError(t, err)
		assert.True(t, acquired)
		// cleanup
		defer func() {
			require.NoError(t, fileLock.Release())
		}()

		targetBaseDir := t.TempDir()
		deployed, err := DeployOneAgent(logger, sourceBaseDir, targetBaseDir, workBaseDir, allTechValue)
		require.NoError(t, err)
		require.False(t, deployed)

		// verify log message that deployment was skipped
		tests.RequireLogMessage(t, logsObserver, "Another instance holds the deployment lock, skipping deployment")

		// verify deployment was NOT performed
		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.Equal(t, NotDeployed, result.Status)
	})

	t.Run("Skip deployment when OneAgent is already deployed", func(t *testing.T) {
		logger, logsObserver := tests.NewTestLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		// setup target directory as already deployed
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, agentVersion)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.Equal(t, Deployed, result.Status)

		workBaseDir := t.TempDir()
		deployed, err := DeployOneAgent(logger, sourceBaseDir, targetBaseDir, workBaseDir, allTechValue)
		require.NoError(t, err)
		require.False(t, deployed)

		// verify DeployOneAgent function detected the deployment was already done
		tests.RequireLogMessage(t, logsObserver, "OneAgent is already deployed")

		// verify the lock file is removed
		lockFilePath := getPathToDeploymentLockFile(workBaseDir)
		_, err = os.Stat(lockFilePath)
		assert.True(t, os.IsNotExist(err), "lock file should be removed after the deployment")
	})

	t.Run("DeployOneAgent returns error when work base folder cannot be created", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.Equal(t, NotDeployed, result.Status)

		// create a parent directory and set its permission to read-only which cause failure in creating work base folder
		workBaseParentDir := t.TempDir()
		err := os.Chmod(workBaseParentDir, dirPerm000)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, os.Chmod(workBaseParentDir, dirPerm700)) // restore permissions on exit to allow cleanup of the temporary directory
		}()

		workBaseDir := filepath.Join(workBaseParentDir, "baseDir")
		deployed, err := DeployOneAgent(logger, sourceBaseDir, targetBaseDir, workBaseDir, allTechValue)
		require.Error(t, err)
		require.False(t, deployed)
		require.Contains(t, err.Error(), "error creating work base folder")
	})

	t.Run("Only one Bootstrapper instance deploy OneAgent in a concurrent environment", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		const (
			agentVersion  = "1.327.30.20251107-111521"
			numGoroutines = 1500
		)

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		workBaseDir := t.TempDir()

		// use a barrier to ensure all goroutines start at the same time
		startBarrier := make(chan struct{})

		var (
			numDeployments int32
			numErrors      int32
		)

		var wg sync.WaitGroup
		wg.Add(numGoroutines)

		for range numGoroutines {
			go func() {
				defer wg.Done()

				<-startBarrier

				deployed, err := DeployOneAgent(logger, sourceBaseDir, targetBaseDir, workBaseDir, allTechValue)
				if err != nil {
					atomic.AddInt32(&numErrors, 1)

					return
				}

				if deployed {
					atomic.AddInt32(&numDeployments, 1)
				}
			}()
		}

		// release all deployment goroutines at the same time
		close(startBarrier)
		wg.Wait()

		// verify that no errors occurred
		assert.Equal(t, int32(0), numErrors)

		// verify that OneAgent is deployed
		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.Equal(t, Deployed, result.Status)

		// verify that only one deployment took place
		assert.Equal(t, int32(1), numDeployments)

		// verify the lock file is removed
		lockFilePath := getPathToDeploymentLockFile(workBaseDir)
		_, err := os.Stat(lockFilePath)
		assert.True(t, os.IsNotExist(err), "lock file should be removed after the deployment")
	})

	t.Run("The stale lock file is removed and OneAgent is successfully deployed", func(t *testing.T) {
		logger, logsObserver := tests.NewTestLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		workBaseDir := t.TempDir()
		lockFilePath := getPathToDeploymentLockFile(workBaseDir)

		// create the stale lock file
		f, err := os.Create(lockFilePath)
		require.NoError(t, err)
		f.Close()

		// set the modification time 30 minutes age to simulate the stale lock scenario
		staleTimestamp := time.Now().Add(-30 * time.Minute)
		err = os.Chtimes(lockFilePath, staleTimestamp, staleTimestamp)
		require.NoError(t, err)

		// the deployment should remove the stale lock file and proceed with the deployment
		targetBaseDir := t.TempDir()
		deployed, err := DeployOneAgent(logger, sourceBaseDir, targetBaseDir, workBaseDir, allTechValue)
		require.NoError(t, err)
		require.True(t, deployed)

		// verify that the stale lock file was detected and removed
		tests.RequireLogMessage(t, logsObserver, "Detected stale lock file, removing it")

		// verify that OneAgent is deployed
		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.NoError(t, result.Error)
		require.Equal(t, Deployed, result.Status)
		require.Equal(t, agentVersion, result.AgentVersion)
	})

	t.Run("Second deployment upgrades OneAgent version", func(t *testing.T) {
		const (
			agentVersion1 = "1.325.22.20251002-101422"
			agentVersion2 = "1.327.30.20251107-111521"
		)

		deployAgentVersions(t, agentVersion1, agentVersion2)
	})

	t.Run("Second deployment downgrades OneAgent version", func(t *testing.T) {
		const (
			agentVersion1 = "1.327.30.20251107-111521"
			agentVersion2 = "1.325.22.20251002-101422"
		)

		deployAgentVersions(t, agentVersion1, agentVersion2)
	})
}

// deployAgentVersions deploys two given OneAgent versions sequentially and verifies the deployment status after each deployment.
func deployAgentVersions(t *testing.T, agentVersion1 string, agentVersion2 string) {
	t.Helper()

	logger, _ := tests.NewTestLogger()

	// setup source directory for OneAgent v1
	sourceAgentV1BaseDir := t.TempDir()
	tests.SetupSourceDirectory(t, sourceAgentV1BaseDir, agentVersion1)

	targetBaseDir := t.TempDir()
	workBaseDir := t.TempDir()

	// check that OneAgent v1 is not deployed
	result := CheckAgentDeploymentStatus(sourceAgentV1BaseDir, targetBaseDir)
	require.Equal(t, NotDeployed, result.Status)

	// deploy OneAgent v1
	deployed, err := DeployOneAgent(logger, sourceAgentV1BaseDir, targetBaseDir, workBaseDir, allTechValue)
	require.NoError(t, err)
	require.True(t, deployed)

	// verify that OneAgent v1 is deployed
	result = CheckAgentDeploymentStatus(sourceAgentV1BaseDir, targetBaseDir)
	require.NoError(t, result.Error)
	require.Equal(t, Deployed, result.Status)
	require.Equal(t, result.AgentVersion, agentVersion1)

	// setup source directory for OneAgent v2
	sourceAgentV2BaseDir := t.TempDir()
	tests.SetupSourceDirectory(t, sourceAgentV2BaseDir, agentVersion2)

	// check that OneAgent v2 is not deployed
	result = CheckAgentDeploymentStatus(sourceAgentV2BaseDir, targetBaseDir)
	require.Equal(t, NotDeployed, result.Status)

	// deploy OneAgent v2
	deployed, err = DeployOneAgent(logger, sourceAgentV2BaseDir, targetBaseDir, workBaseDir, allTechValue)
	require.NoError(t, err)
	require.True(t, deployed)

	// verify that OneAgent v2 is deployed
	result = CheckAgentDeploymentStatus(sourceAgentV2BaseDir, targetBaseDir)
	require.NoError(t, result.Error)
	require.Equal(t, Deployed, result.Status)
	require.Equal(t, result.AgentVersion, agentVersion2)
}
