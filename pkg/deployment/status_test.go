package deployment

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/tests"
	"github.com/stretchr/testify/require"
)

func TestDeploymentStatus(t *testing.T) {
	t.Run("the deployment status is 'Not Deployed' (target is empty)", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, result.Error)
		require.Equal(t, NotDeployed, result.Status)
	})

	t.Run("the deployment status is 'Not Deployed' (target contains another agent version)", func(t *testing.T) {
		const (
			sourceAgentVersion = "1.327.30.20251107-111521"
			targetAgentVersion = "1.325.51.20251103-195814"
		)

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, sourceAgentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, targetAgentVersion, targetAgentVersion)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, result.Error)
		require.Equal(t, NotDeployed, result.Status)
	})

	t.Run("the deployment status is 'Link Missing' (the `active` symlink is missing)", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, "")

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, result.Error)
		require.Equal(t, LinkMissing, result.Status)
	})

	t.Run("the deployment status is 'Link Missing' (the `active` symlink points to another agent version)", func(t *testing.T) {
		const (
			sourceAgentVersion = "1.327.30.20251107-111521"
			targetAgentVersion = "1.325.51.20251103-195814"
		)

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, sourceAgentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, sourceAgentVersion, targetAgentVersion)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, result.Error)
		require.Equal(t, LinkMissing, result.Status)
	})

	t.Run("the deployment status is 'Deployed'", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, agentVersion)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, result.Error)
		require.Equal(t, Deployed, result.Status)
	})

	t.Run("the deployment status is 'Unknown' (due to the target folder permission issue)", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, agentVersion)

		// change the permission of the target folder, so the deployment check must result into an error
		agentTargetDirectory := GetAgentFolder(targetBaseDir, agentVersion)
		parentDir := filepath.Dir(agentTargetDirectory)
		err := os.Chmod(parentDir, 0000)

		defer func() {
			require.NoError(t, os.Chmod(parentDir, 0700)) // restore permissions on exit to allow cleanup of the temporary directory
		}()

		require.NoError(t, err)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.ErrorIs(t, result.Error, syscall.EACCES)
		require.Equal(t, Unknown, result.Status)

		expectedLog := `cannot obtain OneAgent directory info: stat .+: permission denied`
		require.Regexp(t, expectedLog, result.Error.Error())
	})

	t.Run("the deployment status is 'Unknown' (the `active` symlink is a directory)", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, "")

		// create the `active` directory instead of the symlink
		activeDirectoryPath := filepath.Join(targetBaseDir, ActiveLinkPath)
		err := os.MkdirAll(activeDirectoryPath, 0755)
		require.NoError(t, err)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.ErrorContains(t, result.Error, "OneAgent `active` is not a symlink: drwxr-xr-x")
		require.Equal(t, Unknown, result.Status)
	})

	t.Run("the deployment status is 'Unknown' (the oneagent directory is of a file type)", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, agentVersion)

		// remove the agent folder and create the oneagent file instead
		agentTargetDirectory := GetAgentFolder(targetBaseDir, agentVersion)
		err := os.RemoveAll(agentTargetDirectory)
		require.NoError(t, err)

		file, err := os.OpenFile(agentTargetDirectory, os.O_CREATE, dirPerm755)
		require.NoError(t, err)

		defer func() {
			require.NoError(t, file.Close())
		}()

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.ErrorContains(t, result.Error, "OneAgent deployment target is not a directory")
		require.Equal(t, Unknown, result.Status)
	})

	t.Run("the deployment status is 'Unknown' (the installer.version file is not found in the source directory)", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, agentVersion)

		// remove the installer.version file
		versionFilePath := filepath.Join(sourceBaseDir, InstallerVersionFilePath)
		err := os.Remove(versionFilePath)
		require.NoError(t, err)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.ErrorContains(t, result.Error, "failed to determine OneAgent version to deploy")
		require.ErrorIs(t, result.Error, syscall.ENOENT)
		require.Equal(t, Unknown, result.Status)
	})
}
