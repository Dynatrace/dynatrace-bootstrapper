package deployment

import (
	"os"
	"path/filepath"
	"regexp"
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
		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, err)
		require.Equal(t, status, NotDeployed)
	})

	t.Run("the deployment status is 'Not Deployed' (target contains another agent version)", func(t *testing.T) {
		const sourceAgentVersion = "1.327.30.20251107-111521"
		const targetAgentVersion = "1.325.51.20251103-195814"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, sourceAgentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, targetAgentVersion, targetAgentVersion)

		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, err)
		require.Equal(t, status, NotDeployed)
	})

	t.Run("the deployment status is 'Link Missing' (the `active` symlink is missing)", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, "")

		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, err)
		require.Equal(t, status, LinkMissing)
	})

	t.Run("the deployment status is 'Link Missing' (the `active` symlink points to another agent version)", func(t *testing.T) {
		const sourceAgentVersion = "1.327.30.20251107-111521"
		const targetAgentVersion = "1.325.51.20251103-195814"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, sourceAgentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, sourceAgentVersion, targetAgentVersion)

		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, err)
		require.Equal(t, status, LinkMissing)
	})

	t.Run("the deployment status is 'Deployed'", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, agentVersion)

		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.NoError(t, err)
		require.Equal(t, status, Deployed)
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

		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		require.ErrorIs(t, err, syscall.EACCES)
		require.Equal(t, status, Unknown)

		expectedLog := `cannot obtain OneAgent directory info: stat .+: permission denied`
		require.Regexp(t, regexp.MustCompile(expectedLog), err.Error())
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

		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.ErrorContains(t, err, "OneAgent `active` is not a symlink: drwxr-xr-x")
		require.Equal(t, status, Unknown)
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

		file, err := os.OpenFile(agentTargetDirectory, os.O_CREATE, 0755)
		defer func() {
			require.NoError(t, file.Close())
		}()

		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.ErrorContains(t, err, "OneAgent deployment target is not a directory")
		require.Equal(t, status, Unknown)
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

		status, err := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)

		require.ErrorContains(t, err, "failed to determine OneAgent version to deploy")
		require.ErrorIs(t, err, syscall.ENOENT)
		require.Equal(t, status, Unknown)
	})

}
