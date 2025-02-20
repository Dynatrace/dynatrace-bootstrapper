package preload

import (
	"path/filepath"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestConfigure(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		configDir := "path/conf"
		installPath := "/path/install"
		expectedContent := filepath.Join(installPath, libAgentProcPath)

		err := Configure(testLog, fs, configDir, installPath)
		require.NoError(t, err)

		content, err := fs.ReadFile(filepath.Join(configDir, configPath))
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))

	})
}
