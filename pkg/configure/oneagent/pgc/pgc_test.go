package pgc

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
	t.Run("success", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		inputDir := filepath.Join(baseTempDir, "input")
		containerConfigDir := filepath.Join(baseTempDir, "config")

		testData := "test-pgc-data"
		require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, InputFileName), testData))

		err := Configure(testLog, inputDir, containerConfigDir)
		require.NoError(t, err)

		content, err := os.ReadFile(GetDestinationFilePath(containerConfigDir))
		require.NoError(t, err)
		assert.Equal(t, testData, string(content))
	})

	t.Run("missing file == skip", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		inputDir := filepath.Join(baseTempDir, "input")
		containerConfigDir := filepath.Join(baseTempDir, "config")

		err := Configure(testLog, inputDir, containerConfigDir)
		require.NoError(t, err)

		_, err = os.ReadFile(GetDestinationFilePath(containerConfigDir))
		require.True(t, os.IsNotExist(err))
	})
}
