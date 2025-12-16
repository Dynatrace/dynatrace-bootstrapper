package deployment

import (
	"os"
	"regexp"
	"syscall"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/tests"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLogger = zapr.NewLogger(zap.NewExample())

func TestCreateActiveSymlink(t *testing.T) {
	t.Run("the `active` symlink successfully created", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		workDir := t.TempDir()
		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, "")

		agentTargetPath := GetAgentFolder(targetBaseDir, agentVersion)
		err := CreateActiveSymlinkAtomically(testLogger, workDir, agentTargetPath)
		require.NoError(t, err)

		activeSymlink := getPathToActiveLink(agentTargetPath)
		symlinkTarget, err := os.Readlink(activeSymlink)
		require.NoError(t, err)
		require.Equal(t, symlinkTarget, agentVersion)
	})

	t.Run("the existing `active` symlink is successfully updated", func(t *testing.T) {
		const existingAgentVersion = "1.325.51.20251103-195814"
		const targetAgentVersion = "1.327.30.20251107-111521"

		workDir := t.TempDir()
		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, existingAgentVersion, existingAgentVersion)

		agentTargetPath := GetAgentFolder(targetBaseDir, targetAgentVersion)
		err := CreateActiveSymlinkAtomically(testLogger, workDir, agentTargetPath)
		require.NoError(t, err)

		activeSymlink := getPathToActiveLink(agentTargetPath)
		symlinkTarget, err := os.Readlink(activeSymlink)
		require.NoError(t, err)
		require.Equal(t, symlinkTarget, targetAgentVersion)
	})

	t.Run("fail if the target folder is missing", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		workDir := t.TempDir()
		targetBaseDir := t.TempDir()

		agentTargetPath := GetAgentFolder(targetBaseDir, agentVersion)
		err := CreateActiveSymlinkAtomically(testLogger, workDir, agentTargetPath)
		require.ErrorIs(t, err, syscall.ENOENT)

		expectedLog := `failed to rename the temporary symlink: rename .+: no such file or directory`
		require.Regexp(t, regexp.MustCompile(expectedLog), err.Error())
	})

	t.Run("the temporary symlink is removed if the rename operation fails", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		workDir := t.TempDir()
		targetBaseDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetBaseDir, agentVersion, "")

		// change the target folder's permissions to cause the rename of the temporary symlink to fail
		err := os.Chmod(targetBaseDir, 0000)
		defer func() {
			// restore permissions on exit to allow cleanup of the temporary directory
			require.NoError(t, os.Chmod(targetBaseDir, 0700))
		}()
		require.NoError(t, err)

		agentTargetPath := GetAgentFolder(targetBaseDir, agentVersion)
		err = CreateActiveSymlinkAtomically(testLogger, workDir, agentTargetPath)
		require.ErrorIs(t, err, syscall.EACCES)

		entries, err := os.ReadDir(workDir)
		require.NoError(t, err)
		require.Len(t, entries, 0, "the work directory should be empty if the temporary symlink was removed")
	})
}
