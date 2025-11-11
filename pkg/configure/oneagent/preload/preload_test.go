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
}
