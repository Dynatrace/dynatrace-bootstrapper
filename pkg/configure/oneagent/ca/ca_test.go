package ca

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
	configDir := "path/conf"
	inputDir := "/path/input"
	expectedTrusted := "trusted-cert"
	expectedAG := "ag-cert"

	t.Run("success - both present", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		setupTrusted(t, fs, inputDir, expectedTrusted)
		setupAG(t, fs, inputDir, expectedAG)

		err := Configure(testLog, fs, inputDir, configDir)
		require.NoError(t, err)

		certFilePath := filepath.Join(configDir, configBasePath, certsFileName)
		content, err := fs.ReadFile(certFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedTrusted)
		assert.Contains(t, string(content), expectedAG)

		proxyCertFilePath := filepath.Join(configDir, configBasePath, proxyCertsFileName)
		content, err = fs.ReadFile(proxyCertFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedTrusted)
		assert.NotContains(t, string(content), expectedAG)
	})

	t.Run("success - only trusted present", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		setupTrusted(t, fs, inputDir, expectedTrusted)

		err := Configure(testLog, fs, inputDir, configDir)
		require.NoError(t, err)

		certFilePath := filepath.Join(configDir, configBasePath, certsFileName)
		content, err := fs.ReadFile(certFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedTrusted)
		assert.NotContains(t, string(content), expectedAG)

		proxyCertFilePath := filepath.Join(configDir, configBasePath, proxyCertsFileName)
		content, err = fs.ReadFile(proxyCertFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedTrusted)
	})

	t.Run("success - only ag present", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		setupAG(t, fs, inputDir, expectedAG)

		err := Configure(testLog, fs, inputDir, configDir)
		require.NoError(t, err)

		certFilePath := filepath.Join(configDir, configBasePath, certsFileName)
		content, err := fs.ReadFile(certFilePath)
		require.NoError(t, err)
		assert.NotContains(t, string(content), expectedTrusted)
		assert.Contains(t, string(content), expectedAG)

		proxyCertFilePath := filepath.Join(configDir, configBasePath, proxyCertsFileName)
		_, err = fs.ReadFile(proxyCertFilePath)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("missing files == skip", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		err := Configure(testLog, fs, inputDir, configDir)
		require.NoError(t, err)

		certFilePath := filepath.Join(configDir, configBasePath, certsFileName)
		_, err = fs.ReadFile(certFilePath)
		require.True(t, os.IsNotExist(err))

		proxyCertFilePath := filepath.Join(configDir, configBasePath, proxyCertsFileName)
		_, err = fs.ReadFile(proxyCertFilePath)
		require.True(t, os.IsNotExist(err))
	})
}

func setupTrusted(t *testing.T, fs afero.Afero, inputDir, value string) {
	t.Helper()

	require.NoError(t, fsutils.CreateFile(fs, filepath.Join(inputDir, TrustedCertsInputFile), value))
}

func setupAG(t *testing.T, fs afero.Afero, inputDir, value string) {
	t.Helper()

	require.NoError(t, fsutils.CreateFile(fs, filepath.Join(inputDir, AgCertsInputFile), value))
}
