package curl

import (
	"os"
	"path/filepath"
	"testing"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestConfigure(t *testing.T) {
	expectedValue := "123"

	t.Run("success", func(t *testing.T) {
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, "path", "conf")
		inputDir := filepath.Join(tmpDir, "path", "input")

		setupFs(t, inputDir, expectedValue)

		err := Configure(testLog, inputDir, configDir)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(configDir, ConfigPath))
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedValue)
	})

	t.Run("missing file == skip", func(t *testing.T) {
		tmpDir := t.TempDir()
		configDir := filepath.Join(tmpDir, "path", "conf")
		inputDir := filepath.Join(tmpDir, "path", "input")

		err := Configure(testLog, inputDir, configDir)
		require.NoError(t, err)

		_, err = os.ReadFile(filepath.Join(configDir, ConfigPath))
		require.True(t, os.IsNotExist(err))
	})
}

func setupFs(t *testing.T, inputDir, value string) {
	t.Helper()

	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, InputFileName), value))
}
