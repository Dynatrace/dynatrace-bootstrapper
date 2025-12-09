package tests

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

// SetupSourceDirectory creates a mock source directory for testing purposes
func SetupSourceDirectory(t *testing.T, sourceBaseDir, agentVersion string) {
	t.Helper()

	versionFilePath := filepath.Join(sourceBaseDir, "agent/installer.version")
	require.NoError(t, os.MkdirAll(filepath.Dir(versionFilePath), os.ModePerm))
	require.NoError(t, os.WriteFile(versionFilePath, []byte(agentVersion), 0600))

	const agentBinDir = "agent/bin/"
	agentBinFolder := filepath.Join(sourceBaseDir, agentBinDir, agentVersion)
	require.NoError(t, os.MkdirAll(agentBinFolder, 0700))
}

// SetupTargetDirectory creates a mock target directory for testing purposes.
//
// Parameters:
// - targetBaseDir: the base path to the target directory (e.g., "/home/dynatrace/")
// - agentVersionDir: the name of the versioned OneAgent directory (e.g., "1.325.51.20251103-195814")
// - activeLinkAgentVersion: the OneAgent version that the `active` symbolic link will point to (e.g., "1.323.49.20251103-195824")
//
// Note:
// The values of `agentVersionDir` and `activeLinkAgentVersion` parameters may differ depending on a unit test scenario.
func SetupTargetDirectory(t *testing.T, targetBaseDir, agentVersionDir, activeLinkAgentVersion string) {
	t.Helper()

	oneAgentDirPath := filepath.Join(targetBaseDir, "oneagent", agentVersionDir)
	err := os.MkdirAll(oneAgentDirPath, 0755)
	require.NoError(t, err)

	if activeLinkAgentVersion != "" {
		activeLinkPath := filepath.Join(targetBaseDir, "oneagent/active")
		err = os.Symlink(activeLinkAgentVersion, activeLinkPath)
		require.NoError(t, err)
	}
}
