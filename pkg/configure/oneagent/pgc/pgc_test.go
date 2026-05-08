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
		targetDir := filepath.Join(baseTempDir, "target")

		testData := "test-pgc-data"
		require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, InputFileName), testData))

		err := Configure(testLog, inputDir, targetDir)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(targetDir, DestinationPath))
		require.NoError(t, err)
		assert.Equal(t, testData, string(content))
	})

	t.Run("missing file == skip", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		inputDir := filepath.Join(baseTempDir, "input")
		targetDir := filepath.Join(baseTempDir, "target")

		err := Configure(testLog, inputDir, targetDir)
		require.NoError(t, err)

		_, err = os.ReadFile(filepath.Join(targetDir, DestinationPath))
		require.True(t, os.IsNotExist(err))
	})
}
