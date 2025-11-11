package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/stretchr/testify/require"
)

func TestBootstrapper(t *testing.T) {
	t.Run("should validate required flags - missing flags -> error", func(t *testing.T) {
		cmd := New()

		err := cmd.Execute()

		require.Error(t, err)
	})
	t.Run("should validate required flags - present flags -> no error", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupSource(t, tmpDir, "123")

		cmd := New()
		cmd.SetArgs([]string{"--source", tmpDir, "--target", t.TempDir()})

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("--suppress-error=true -> no error", func(t *testing.T) {
		cmd := New()
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\"})

		err := cmd.Execute()

		require.Error(t, err)

		cmd = New()
		// Note: we can't skip/mask the validation of the required flags
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\", "--suppress-error"})

		err = cmd.Execute()

		require.NoError(t, err)

		cmd = New()
		// Note: we can't skip/mask the validation of the required flags
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\", "--suppress-error", "true"})

		err = cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("should allow unknown flags -> no error", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupSource(t, tmpDir, "123")

		cmd := New()
		cmd.SetArgs([]string{"--source", tmpDir, "--target", t.TempDir(), "--unknown", "--flag", "value"})

		err := cmd.Execute()

		require.NoError(t, err)
	})
}

func setupSource(t *testing.T, folder, version string) {
	t.Helper()

	versionFilePath := filepath.Join(folder, move.InstallerVersionFilePath)
	require.NoError(t, os.MkdirAll(filepath.Dir(versionFilePath), os.ModePerm))
	require.NoError(t, os.WriteFile(versionFilePath, []byte(version), 0600))

	agentBinFolder := filepath.Join(folder, filepath.Dir(move.CurrentDir), version)
	require.NoError(t, os.MkdirAll(agentBinFolder, 0700))
}
