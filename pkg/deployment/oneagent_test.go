package deployment

import (
	"os"
	"path/filepath"
	"regexp"
	"syscall"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/tests"
	"github.com/stretchr/testify/require"
)

func TestCopyAgent(t *testing.T) {
	t.Run("Successful copy from Source to Target", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"
		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		workBaseDir := t.TempDir()
		targetBaseDir := t.TempDir()
		agentFolder := GetAgentFolder(targetBaseDir, agentVersion)
		err := CopyAgent(testLogger, sourceBaseDir, agentFolder, workBaseDir, allTechValue)
		require.NoError(t, err)

		result := CheckAgentDeploymentStatus(sourceBaseDir, targetBaseDir)
		// the `active` symlink must be missing because it is created after the copy operation
		require.Equal(t, result.Status, LinkMissing)
		require.Equal(t, result.AgentVersion, agentVersion)

		// the `current` symlink must exist because it is created during the copy operation
		currentSymlinkPath := filepath.Join(agentFolder, move.CurrentDir)
		stat, err := os.Lstat(currentSymlinkPath)
		require.NoError(t, err)
		require.True(t, stat.Mode()&os.ModeSymlink != 0, "current should be a symlink type")
	})

	t.Run("Cannot copy due to a permission error in the target directory", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceBaseDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceBaseDir, agentVersion)

		targetBaseDir := t.TempDir()
		err := os.Chmod(targetBaseDir, 0000)
		defer func() {
			require.NoError(t, os.Chmod(targetBaseDir, 0755)) // restore permissions on exit to allow cleanup of the target directory
		}()

		workBaseDir := t.TempDir()
		agentFolder := GetAgentFolder(targetBaseDir, agentVersion)
		err = CopyAgent(testLogger, sourceBaseDir, agentFolder, workBaseDir, allTechValue)
		require.ErrorIs(t, err, syscall.EACCES)

		expectedLog := `failed to create the target folder: mkdir .+: permission denied`
		require.Regexp(t, regexp.MustCompile(expectedLog), err.Error())
	})

}
