package preload

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestConfigure(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "conf")
		installPath := filepath.Join(baseTempDir, "install")
		expectedContent := filepath.Join(installPath, LibAgentProcPath)

		err := Configure(testLog, configDir, installPath)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(configDir, ConfigPath))
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	})

	t.Run("relative install path is rejected", func(t *testing.T) {
		err := Configure(testLog, t.TempDir(), "relative/path")
		require.Error(t, err)
	})

	t.Run("install path with separator is rejected", func(t *testing.T) {
		for _, path := range []string{
			"/valid/path,/injected/path",
			"/valid/path:/other",
			"/valid/path\x00/other",
		} {
			require.Error(t, Configure(testLog, t.TempDir(), path))
		}
	})

	t.Run("install path with whitespace is rejected", func(t *testing.T) {
		for _, path := range []string{
			"/valid/path\n/injected/path",
			"/valid/path\r/other",
			"/valid/path\t/other",
			"/valid/path\x00/other",
		} {
			require.Error(t, Configure(testLog, t.TempDir(), path))
		}
	})

	t.Run("unclean install path is rejected", func(t *testing.T) {
		err := Configure(testLog, t.TempDir(), "/valid/../path")
		require.Error(t, err)
	})

	// Valid but rejected: spaces in paths are indistinguishable from
	// whitespace separators in ld.so.preload, so we treat them as invalid.
	t.Run("path with space is rejected despite being a valid linux path", func(t *testing.T) {
		err := Configure(testLog, t.TempDir(), "/opt/my agent/dynatrace")
		require.Error(t, err)
	})
}
