package curl

import (
	"os"
	"path/filepath"
	"testing"

	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/zapr"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestConfigure(t *testing.T) {
	expectedValue := "123"
	configDir := "path/conf"
	inputDir := "/path/input"

	t.Run("success", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		setupFs(t, fs, inputDir, expectedValue)

		err := Configure(testLog, fs, inputDir, configDir)
		require.NoError(t, err)

		content, err := fs.ReadFile(filepath.Join(configDir, ConfigPath))
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedValue)
	})

	t.Run("missing file == skip", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		err := Configure(testLog, fs, inputDir, configDir)
		require.NoError(t, err)

		_, err = fs.ReadFile(filepath.Join(configDir, ConfigPath))
		require.True(t, os.IsNotExist(err))
	})
}

func setupFs(t *testing.T, fs afero.Afero, inputDir, value string) {
	t.Helper()

	require.NoError(t, fsutils.CreateFile(fs, filepath.Join(inputDir, InputFileName), value))
}
