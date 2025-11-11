package pmc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
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
		configDir := filepath.Join(baseTempDir, "config", "container")
		installPath := filepath.Join(baseTempDir, "install")
		targetDir := filepath.Join(baseTempDir, "target")
		inputDir := filepath.Join(baseTempDir, "input")

		source := ruxit.ProcConf{
			Properties: []ruxit.Property{
				{
					Section: "test",
					Key:     "key",
					Value:   "value",
				},
				{
					Section: "test",
					Key:     "source",
					Value:   "source",
				},
			},
		}

		override := ruxit.ProcConf{
			Properties: []ruxit.Property{
				{
					Section: "test",
					Key:     "key",
					Value:   "override",
				},
				{
					Section: "test",
					Key:     "add",
					Value:   "add",
				},
			},
			InstallPath: &installPath,
		}

		setupInputFs(t, inputDir, override)
		setupTargetFs(t, targetDir, source)

		err := Configure(testLog, inputDir, targetDir, configDir, installPath)
		require.NoError(t, err)

		content, err := os.ReadFile(GetSourceRuxitAgentProcFilePath(targetDir))
		require.NoError(t, err)
		assert.Equal(t, source.ToString(), string(content))

		content, err = os.ReadFile(GetDestinationRuxitAgentProcFilePath(configDir))
		require.NoError(t, err)
		assert.Equal(t, source.Merge(override).ToString(), string(content))
	})
	t.Run("missing file == skip", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config", "container")
		installPath := filepath.Join(baseTempDir, "install")
		targetDir := filepath.Join(baseTempDir, "target")
		inputDir := filepath.Join(baseTempDir, "input")

		source := ruxit.ProcConf{
			Properties: []ruxit.Property{
				{
					Section: "test",
					Key:     "key",
					Value:   "value",
				},
				{
					Section: "test",
					Key:     "source",
					Value:   "source",
				},
			},
		}

		setupTargetFs(t, targetDir, source)

		err := Configure(testLog, inputDir, targetDir, configDir, installPath)
		require.NoError(t, err)

		content, err := os.ReadFile(GetSourceRuxitAgentProcFilePath(targetDir))
		require.NoError(t, err)
		assert.Equal(t, source.ToString(), string(content))

		_, err = os.ReadFile(GetDestinationRuxitAgentProcFilePath(configDir))
		require.True(t, os.IsNotExist(err))
	})
}

func setupInputFs(t *testing.T, inputDir string, value ruxit.ProcConf) {
	t.Helper()

	rawValue, err := json.Marshal(value)
	require.NoError(t, err)
	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, InputFileName), string(rawValue)))
}

func setupTargetFs(t *testing.T, targetDir string, value ruxit.ProcConf) {
	t.Helper()

	require.NoError(t, fsutils.CreateFile(filepath.Join(targetDir, SourceRuxitAgentProcPath), value.ToString()))
}
