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
	t.Run("simple copy", func(t *testing.T) {
		tmpDir := t.TempDir()
		sourceDir := filepath.Join(tmpDir, "source")
		targetDir := filepath.Join(tmpDir, "target")

		// Create source directory and files
		workDir := filepath.Join(tmpDir, "work")
		_ = os.MkdirAll(sourceDir, 0755)
		_ = os.WriteFile(sourceDir+"/file1.txt", []byte("file1 content"), 0600)
		_ = os.WriteFile(sourceDir+"/file2.txt", []byte("file2 content"), 0600)
		setupSource(t, sourceDir, "123")

		workFolder = workDir

		technology = " " + AllTechValue + " "

		err := Execute(testLog, sourceDir, targetDir)
		require.NoError(t, err)

		// Check if the target directory and files exist
		_, err = os.Stat(targetDir)
		require.NoError(t, err)

		_, err = os.Stat(targetDir + "/file1.txt")
		require.NoError(t, err)

		_, err = os.Stat(targetDir + "/file2.txt")
		require.NoError(t, err)

		// Check the content of the copied files
		content, err := os.ReadFile(targetDir + "/file1.txt")
		require.NoError(t, err)
		assert.Equal(t, "file1 content", string(content))

		content, err = os.ReadFile(targetDir + "/file2.txt")
		require.NoError(t, err)
		assert.Equal(t, "file2 content", string(content))

		// Check the cleanup happened of the copied files
		_, err = os.Stat(workDir)
		require.Error(t, err)
	})
	t.Run("execute with technology param", func(t *testing.T) {
		tmpDir := t.TempDir()
		sourceDir := filepath.Join(tmpDir, "source")
		targetDir := filepath.Join(tmpDir, "target")

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

		technologyList := "java"
		_ = os.MkdirAll(sourceDir, 0755)
		_ = os.WriteFile(sourceDir+"/manifest.json", []byte(manifestContent), 0600)
		_ = os.WriteFile(sourceDir+"/fileA1.txt", []byte("fileA1 content"), 0600)
		_ = os.WriteFile(sourceDir+"/fileA2.txt", []byte("fileA2 content"), 0600)
		setupSource(t, sourceDir, "123")

		technology = technologyList

		err := Execute(testLog, sourceDir, targetDir)
		require.NoError(t, err)

		// Check if the target directory and files exist
		_, err = os.Stat(targetDir)
		require.NoError(t, err)

		_, err = os.Stat(targetDir + "/fileA1.txt")
		require.NoError(t, err)

		_, err = os.Stat(targetDir + "/fileA2.txt")
		require.Error(t, err)

		// Check the content of the copied files
		content, err := os.ReadFile(targetDir + "/fileA1.txt")
		require.NoError(t, err)
		assert.Equal(t, "fileA1 content", string(content))

		content, err = os.ReadFile(targetDir + "/fileA2.txt")
		require.Error(t, err)
		assert.Empty(t, string(content))
	})
}

func setupSource(t *testing.T, folder, version string) {
	t.Helper()

	versionFilePath := filepath.Join(folder, impl.InstallerVersionFilePath)
	require.NoError(t, os.MkdirAll(filepath.Dir(versionFilePath), os.ModePerm))
	require.NoError(t, os.WriteFile(versionFilePath, []byte(version), 0600))

	agentBinFolder := filepath.Join(folder, filepath.Dir(impl.CurrentDir), version)
	require.NoError(t, os.MkdirAll(agentBinFolder, 0700))
}
