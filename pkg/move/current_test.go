package move

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateCurrentSymlink(t *testing.T) {
	expectedVersion := "1.239.14.20220325-164521"

	t.Run("happy path", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupAgentBin(t, tmpDir, expectedVersion)
		setupVersionFile(t, tmpDir, expectedVersion)

		err := CreateCurrentSymlink(testLog, tmpDir)
		require.NoError(t, err)

		linkedDir, err := os.Readlink(filepath.Join(tmpDir, CurrentDir))
		require.NoError(t, err)
		assert.NotEmpty(t, linkedDir)
	})

	t.Run("no fail if current dir already exists", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupAgentBin(t, tmpDir, expectedVersion)
		setupCurrentBin(t, tmpDir)
		setupVersionFile(t, tmpDir, expectedVersion)

		err := CreateCurrentSymlink(testLog, tmpDir)
		require.NoError(t, err)
	})

	t.Run("fail if version file is missing", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupAgentBin(t, tmpDir, expectedVersion)

		err := CreateCurrentSymlink(testLog, tmpDir)
		require.Error(t, err)
	})

	t.Run("fail if bin folder is missing", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupVersionFile(t, tmpDir, expectedVersion)

		err := CreateCurrentSymlink(testLog, tmpDir)
		require.Error(t, err)
	})
}

func setupAgentBin(t *testing.T, folder, version string) {
	t.Helper()

	agentBinFolder := filepath.Join(folder, filepath.Dir(CurrentDir), version)
	require.NoError(t, os.MkdirAll(agentBinFolder, 0700))
}

func setupCurrentBin(t *testing.T, folder string) {
	t.Helper()

	agentCurrentFolder := filepath.Join(folder, CurrentDir)
	require.NoError(t, os.MkdirAll(agentCurrentFolder, 0700))
}

func setupVersionFile(t *testing.T, folder, version string) {
	t.Helper()

	versionFilePath := filepath.Join(folder, InstallerVersionFilePath)
	require.NoError(t, os.MkdirAll(filepath.Dir(versionFilePath), os.ModePerm))
	require.NoError(t, os.WriteFile(versionFilePath, []byte(version), 0600))
}
