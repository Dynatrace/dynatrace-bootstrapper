package ca

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
	expectedTrusted := "trusted-cert"
	expectedAG := "ag-cert"

	t.Run("success - both present", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")
		inputDir := filepath.Join(baseTempDir, "input")

		setupTrusted(t, inputDir, expectedTrusted)
		setupAG(t, inputDir, expectedAG)

		err := Configure(testLog, inputDir, configDir)
		require.NoError(t, err)

		certFilePath := filepath.Join(configDir, ConfigBasePath, CertsFileName)
		content, err := os.ReadFile(certFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedTrusted)
		assert.Contains(t, string(content), expectedAG)

		proxyCertFilePath := filepath.Join(configDir, ConfigBasePath, ProxyCertsFileName)
		content, err = os.ReadFile(proxyCertFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedTrusted)
		assert.NotContains(t, string(content), expectedAG)
	})

	t.Run("success - only trusted present", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")
		inputDir := filepath.Join(baseTempDir, "input")

		setupTrusted(t, inputDir, expectedTrusted)

		err := Configure(testLog, inputDir, configDir)
		require.NoError(t, err)

		certFilePath := filepath.Join(configDir, ConfigBasePath, CertsFileName)
		content, err := os.ReadFile(certFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedTrusted)
		assert.NotContains(t, string(content), expectedAG)

		proxyCertFilePath := filepath.Join(configDir, ConfigBasePath, ProxyCertsFileName)
		content, err = os.ReadFile(proxyCertFilePath)
		require.NoError(t, err)
		assert.Contains(t, string(content), expectedTrusted)
	})

	t.Run("success - only ag present", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")
		inputDir := filepath.Join(baseTempDir, "input")

		setupAG(t, inputDir, expectedAG)

		err := Configure(testLog, inputDir, configDir)
		require.NoError(t, err)

		certFilePath := filepath.Join(configDir, ConfigBasePath, CertsFileName)
		content, err := os.ReadFile(certFilePath)
		require.NoError(t, err)
		assert.NotContains(t, string(content), expectedTrusted)
		assert.Contains(t, string(content), expectedAG)

		proxyCertFilePath := filepath.Join(configDir, ConfigBasePath, ProxyCertsFileName)
		_, err = os.ReadFile(proxyCertFilePath)
		require.True(t, os.IsNotExist(err))
	})

	t.Run("missing files == skip", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")
		inputDir := filepath.Join(baseTempDir, "input")

		err := Configure(testLog, inputDir, configDir)
		require.NoError(t, err)

		certFilePath := filepath.Join(configDir, ConfigBasePath, CertsFileName)
		_, err = os.ReadFile(certFilePath)
		require.True(t, os.IsNotExist(err))

		proxyCertFilePath := filepath.Join(configDir, ConfigBasePath, ProxyCertsFileName)
		_, err = os.ReadFile(proxyCertFilePath)
		require.True(t, os.IsNotExist(err))
	})
}

func setupTrusted(t *testing.T, inputDir, value string) {
	t.Helper()

	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, TrustedCertsInputFile), value))
}

func setupAG(t *testing.T, inputDir, value string) {
	t.Helper()

	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, AgCertsInputFile), value))
}
