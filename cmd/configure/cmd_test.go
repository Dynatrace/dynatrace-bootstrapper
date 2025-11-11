package configure

import (
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/enrichment/endpoint"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/ca"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/curl"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
	fsutils "github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

// Only checking the counts of files in the folders, checking exact paths and contents are done in the sub-package tests.
func TestSetupOneAgent(t *testing.T) {
	podAttributes = []string{
		"k8s.pod.name=pod1",
		"k8s.pod.uid=123",
		"k8s.namespace.name=default",
	}

	containerNames := []string{"test-container-name", "other-container-name"}
	containerAttributes = []string{
		`{"container_image.registry": "some.reg.io", "container_image.repository": "test-repo", "container_image.tags": "latest", "container_image.digest": "sha256:abcd1234", "k8s.container.name": "test-container-name"}`,
		`{"container_image.registry": "some.reg.io", "container_image.repository": "test-repo", "container_image.tags": "latest", "container_image.digest": "sha256:abcd1234", "k8s.container.name": "other-container-name"}`,
	}

	t.Run("success", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir = filepath.Join(baseTempDir, "conf")
		inputDir = filepath.Join(baseTempDir, "input")
		targetFolder := filepath.Join(baseTempDir, "target")

		setupInputFs(t, inputDir)
		setupTargetFs(t, targetFolder)

		preExecuteConfigCount := countFiles(t, configDir)
		require.Equal(t, 0, preExecuteConfigCount)

		preExecuteTargetCount := countFiles(t, targetFolder)
		require.Equal(t, 1, preExecuteTargetCount) // for ruxitagentproc.conf, you need a source file

		err := SetupOneAgent(testLog, targetFolder)
		require.NoError(t, err)

		expectedContainerSpecificConfigCount := 5 // curl(1) + ca(2) + conf(1) + ruxitagentproc.conf(1)

		for _, name := range containerNames {
			containerConfigFolder := filepath.Join(configDir, name)

			containerSpecificConfigCount := countFiles(t, containerConfigFolder)
			require.Equal(t, expectedContainerSpecificConfigCount, containerSpecificConfigCount)
		}

		expectedPostExecuteConfigCount := 1 + len(containerNames)*expectedContainerSpecificConfigCount // preload(1) + len(containers) * container-specific-files
		postExecuteConfigCount := countFiles(t, configDir)
		require.Equal(t, expectedPostExecuteConfigCount, postExecuteConfigCount)

		postExecuteTargetCount := countFiles(t, targetFolder)
		require.Equal(t, preExecuteTargetCount, postExecuteTargetCount) // no change to the target folder during configuration
	})

	t.Run("no input-directory ==> do nothing", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir = filepath.Join(baseTempDir, "conf")
		targetFolder := filepath.Join(baseTempDir, "target")
		inputDir = ""

		err := SetupOneAgent(testLog, targetFolder)
		require.NoError(t, err)

		postExecuteConfigCount := countFiles(t, configDir)
		require.Equal(t, 0, postExecuteConfigCount)

		postExecuteTargetCount := countFiles(t, targetFolder)
		require.Equal(t, 0, postExecuteTargetCount)
	})

	t.Run("no config-directory ==> do nothing", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		inputDir = filepath.Join(baseTempDir, "input")
		targetFolder := filepath.Join(baseTempDir, "target")
		configDir = ""

		err := SetupOneAgent(testLog, targetFolder)
		require.NoError(t, err)

		postExecuteConfigCount := countFiles(t, configDir)
		require.Equal(t, 0, postExecuteConfigCount)

		postExecuteTargetCount := countFiles(t, targetFolder)
		require.Equal(t, 0, postExecuteTargetCount)
	})
}

func TestEnrichWithMetadata(t *testing.T) {
	podAttributes = []string{
		"k8s.pod.name=pod1",
		"k8s.pod.uid=123",
		"k8s.namespace.name=default",
	}

	containerNames := []string{"test-container-name", "other-container-name"}
	containerAttributes = []string{
		`{"container_image.registry": "some.reg.io", "container_image.repository": "test-repo", "container_image.tags": "latest", "container_image.digest": "sha256:abcd1234", "k8s.container.name": "test-container-name"}`,
		`{"container_image.registry": "some.reg.io", "container_image.repository": "test-repo", "container_image.tags": "latest", "container_image.digest": "sha256:abcd1234", "k8s.container.name": "other-container-name"}`,
	}

	t.Run("success", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir = filepath.Join(baseTempDir, "conf")
		inputDir = filepath.Join(baseTempDir, "input")

		setupInputFs(t, inputDir)

		preExecuteConfigCount := countFiles(t, configDir)
		require.Equal(t, 0, preExecuteConfigCount)

		err := EnrichWithMetadata(testLog)
		require.NoError(t, err)

		expectedContainerSpecificConfigCount := 3 // endpoint(1) + metadata(2)

		for _, name := range containerNames {
			containerConfigFolder := filepath.Join(configDir, name)

			containerSpecificConfigCount := countFiles(t, containerConfigFolder)
			require.Equal(t, expectedContainerSpecificConfigCount, containerSpecificConfigCount)
		}

		expectedPostExecuteConfigCount := len(containerNames) * expectedContainerSpecificConfigCount // len(containers) * container-specific-files
		postExecuteConfigCount := countFiles(t, configDir)
		require.Equal(t, expectedPostExecuteConfigCount, postExecuteConfigCount)
	})

	t.Run("no input-directory ==> do nothing", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir = filepath.Join(baseTempDir, "conf")
		inputDir = ""

		err := EnrichWithMetadata(testLog)
		require.NoError(t, err)

		postExecuteConfigCount := countFiles(t, configDir)
		require.Equal(t, 0, postExecuteConfigCount)
	})

	t.Run("no config-directory ==> do nothing", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir = ""
		inputDir = filepath.Join(baseTempDir, "input")

		err := EnrichWithMetadata(testLog)
		require.NoError(t, err)

		postExecuteConfigCount := countFiles(t, configDir)
		require.Equal(t, 0, postExecuteConfigCount)
	})
}

func countFiles(t *testing.T, path string) int {
	t.Helper()

	count := 0
	_ = filepath.Walk(path, func(_ string, info fs.FileInfo, err error) error {
		if os.IsNotExist(err) {
			return nil
		}

		if !info.IsDir() {
			count++
		}

		require.NoError(t, err)

		return nil
	})

	return count
}

func setupInputFs(t *testing.T, inputDir string) {
	t.Helper()

	// endpoint
	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, endpoint.InputFileName), "endpoint"))

	// ca
	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, ca.TrustedCertsInputFile), "trusted"))
	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, ca.AgCertsInputFile), "ag"))

	// curl
	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, curl.InputFileName), "123"))

	// pmc
	procConf := ruxit.ProcConf{
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
	}

	rawProcConf, err := json.Marshal(procConf)
	require.NoError(t, err)
	require.NoError(t, fsutils.CreateFile(filepath.Join(inputDir, pmc.InputFileName), string(rawProcConf)))
}

func setupTargetFs(t *testing.T, targetDir string) {
	t.Helper()

	procConf := ruxit.ProcConf{
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

	require.NoError(t, fsutils.CreateFile(filepath.Join(targetDir, pmc.SourceRuxitAgentProcPath), procConf.ToString()))
}
