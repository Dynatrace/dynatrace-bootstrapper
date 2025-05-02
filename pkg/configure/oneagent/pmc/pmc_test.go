package pmc

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/zapr"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestConfigure(t *testing.T) {
	targetDir := "path/target"
	inputDir := "/path/input"
	configDir := "/path/config/container"
	installPath := "path/install"

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

	t.Run("success", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		setupInputFs(t, fs, inputDir, override)
		setupTargetFs(t, fs, targetDir, source)

		err := Configure(testLog, fs, inputDir, targetDir, configDir, installPath)
		require.NoError(t, err)

		content, err := fs.ReadFile(GetSourceRuxitAgentProcPath(targetDir))
		require.NoError(t, err)
		assert.Equal(t, source.ToString(), string(content))

		content, err = fs.ReadFile(GetDestinationRuxitAgentProcPath(configDir))
		require.NoError(t, err)
		assert.Equal(t, source.Merge(override).ToString(), string(content))
	})

	t.Run("missing file == skip", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		setupTargetFs(t, fs, targetDir, source)

		err := Configure(testLog, fs, inputDir, targetDir, configDir, installPath)
		require.NoError(t, err)

		content, err := fs.ReadFile(GetSourceRuxitAgentProcPath(targetDir))
		require.NoError(t, err)
		assert.Equal(t, source.ToString(), string(content))

		_, err = fs.ReadFile(GetDestinationRuxitAgentProcPath(configDir))
		require.True(t, os.IsNotExist(err))
	})
}

func setupInputFs(t *testing.T, fs afero.Afero, inputDir string, value ruxit.ProcConf) {
	t.Helper()

	rawValue, err := json.Marshal(value)
	require.NoError(t, err)
	require.NoError(t, fsutils.CreateFile(fs, filepath.Join(inputDir, InputFileName), string(rawValue)))
}

func setupTargetFs(t *testing.T, fs afero.Afero, targetDir string, value ruxit.ProcConf) {
	t.Helper()

	require.NoError(t, fsutils.CreateFile(fs, filepath.Join(targetDir, SourceRuxitAgentProcPath), value.ToString()))
}
