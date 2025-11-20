package move

import (
	"os"
	"path/filepath"
	"testing"

	impl "github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestExecute(t *testing.T) {
	const (
		file1 = "fileA1.txt"
		file2 = "fileA2.txt"
	)

	t.Run("simple copy", func(t *testing.T) {
		tmpDir := t.TempDir()
		sourceDir := filepath.Join(tmpDir, "source")
		targetDir := filepath.Join(tmpDir, "target")
		workDir := filepath.Join(tmpDir, "work")

		files := map[string]string{
			file1: "file1 content",
			file2: "file2 content",
		}

		setupSource(t, sourceDir, "123", files)

		workFolder = workDir

		technology = " " + AllTechValue + " "

		err := Execute(testLog, sourceDir, targetDir)
		require.NoError(t, err)

		verifyTarget(t, targetDir, files)

		// Check the cleanup happened of the copied files
		_, err = os.Stat(workDir)
		require.Error(t, err)
	})
	t.Run("execute with technology param", func(t *testing.T) {
		tmpDir := t.TempDir()
		sourceDir := filepath.Join(tmpDir, "source")
		targetDir := filepath.Join(tmpDir, "target")

		manifestFile := "manifest.json"
		manifestContent := `{
			"version": "1.0",
			"technologies": {
				"java": {
					"x86": [
						{"path": "fileA1.txt", "version": "1.0", "md5": "abc123"},
						{"path": "agent/installer.version", "version": "1.0", "md5": "abc123"},
						{"path": "agent/bin/123", "version": "1.0", "md5": "abc123"}
					]
				},
				"python": {
					"arm": [
						{"path": "fileA2.txt", "version": "1.0", "md5": "ghi789"}
					]
				}
			}
		}`

		files := map[string]string{
			manifestFile: manifestContent,
			file1:        "file1 content",
			file2:        "file2 content",
		}

		expectedFiles := map[string]string{
			file1: "file1 content",
		}

		setupSource(t, sourceDir, "123", files)

		technologyList := "java"

		technology = technologyList

		err := Execute(testLog, sourceDir, targetDir)
		require.NoError(t, err)

		verifyTarget(t, targetDir, expectedFiles, file2)
	})
}

func setupSource(t *testing.T, folder, version string, filesToCreate map[string]string) {
	t.Helper()

	require.NoError(t, os.Mkdir(folder, os.ModePerm))

	for path, content := range filesToCreate {
		require.NoError(t, os.WriteFile(filepath.Join(folder, path), []byte(content), 0600))
	}

	versionFilePath := filepath.Join(folder, impl.InstallerVersionFilePath)
	require.NoError(t, os.MkdirAll(filepath.Dir(versionFilePath), os.ModePerm))
	require.NoError(t, os.WriteFile(versionFilePath, []byte(version), 0600))

	agentBinFolder := filepath.Join(folder, filepath.Dir(impl.CurrentDir), version)
	require.NoError(t, os.MkdirAll(agentBinFolder, 0700))
}

func verifyTarget(t *testing.T, folder string, copiedFiles map[string]string, missingFiles ...string) {
	t.Helper()

	for path, expectedContent := range copiedFiles {
		content, err := os.ReadFile(filepath.Join(folder, path))
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	}

	for _, path := range missingFiles {
		assert.NoFileExists(t, filepath.Join(folder, path))
	}
}
