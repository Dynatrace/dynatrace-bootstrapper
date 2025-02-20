package move

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCopyFolderWithTechnologyFiltering(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	sourceDir := "/source"
	targetDir := "/target"

	_ = fs.MkdirAll(sourceDir, 0755)
	_ = fs.MkdirAll(targetDir, 0755)

	manifestContent := `{
        "version": "1.0",
        "technologies": {
            "java": {
                "x86": [
                    {"path": "fileA1.txt", "version": "1.0", "md5": "abc123"},
                    {"path": "fileA2.txt", "version": "1.0", "md5": "def456"}
                ]
            },
            "python": {
                "arm": [
                    {"path": "fileB1.txt", "version": "1.0", "md5": "ghi789"}
                ]
            }
        }
    }`

	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "manifest.json"), []byte(manifestContent), 0644)
	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "fileA1.txt"), []byte("java a1"), 0644)
	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "fileA2.txt"), []byte("java a2"), 0644)
	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "fileB1.txt"), []byte("python b1"), 0644)
	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "fileC1.txt"), []byte("unrelated"), 0644)

	t.Run("copy with single technology filter", func(t *testing.T) {
		t.Cleanup(func() {
			_ = fs.RemoveAll(targetDir)
			_ = fs.MkdirAll(targetDir, 0755)
		})

		technology = "java"
		err := copyByTechnology(testLog, fs, sourceDir, targetDir)
		require.NoError(t, err)

		assertFileExists(t, fs, filepath.Join(targetDir, "fileA1.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileA2.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileB1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileC1.txt"))
	})
	t.Run("copy with multiple technology filter", func(t *testing.T) {
		t.Cleanup(func() {
			_ = fs.RemoveAll(targetDir)
			_ = fs.MkdirAll(targetDir, 0755)
		})

		technology = "java,python"
		err := copyByTechnology(testLog, fs, sourceDir, targetDir)
		require.NoError(t, err)

		assertFileExists(t, fs, filepath.Join(targetDir, "fileA1.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileA2.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileB1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileC1.txt"))
	})
	t.Run("copy with invalid technology filter", func(t *testing.T) {
		t.Cleanup(func() {
			_ = fs.RemoveAll(targetDir)
			_ = fs.MkdirAll(targetDir, 0755)
		})

		technology = "php"
		err := copyByTechnology(testLog, fs, sourceDir, targetDir)
		require.NoError(t, err)

		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileA1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileA2.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileB1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileC1.txt"))
	})
}

func TestFilterFilesByTechnology(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	sourceDir := "/source"
	_ = fs.MkdirAll(sourceDir, 0755)
	manifestContent := `{
        "version": "1.0",
        "technologies": {
            "java": {
                "x86": [
                    {"path": "fileA1.txt", "version": "1.0", "md5": "abc123"},
                    {"path": "fileA2.txt", "version": "1.0", "md5": "def456"}
                ]
            },
            "python": {
                "arm": [
                    {"path": "fileB1.txt", "version": "1.0", "md5": "ghi789"}
                ]
            }
        }
    }`
	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "manifest.json"), []byte(manifestContent), 0644)
	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "fileA1.txt"), []byte("a1 content"), 0644)
	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "fileA2.txt"), []byte("a2 content"), 0644)
	_ = afero.WriteFile(fs, filepath.Join(sourceDir, "fileB1.txt"), []byte("b1 content"), 0644)

	t.Run("filter single technology", func(t *testing.T) {
		paths, err := filterFilesByTechnology(testLog, fs, sourceDir, []string{"java"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			filepath.Join(sourceDir, "fileA1.txt"),
			filepath.Join(sourceDir, "fileA2.txt"),
		}, paths)
	})
	t.Run("filter multiple technologies", func(t *testing.T) {
		paths, err := filterFilesByTechnology(testLog, fs, sourceDir, []string{"java", "python"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			filepath.Join(sourceDir, "fileA1.txt"),
			filepath.Join(sourceDir, "fileA2.txt"),
			filepath.Join(sourceDir, "fileB1.txt"),
		}, paths)
	})
	t.Run("not filter non-existing technology", func(t *testing.T) {
		paths, err := filterFilesByTechnology(testLog, fs, sourceDir, []string{"php"})
		require.NoError(t, err)
		assert.Empty(t, paths)
	})
	t.Run("filter with missing manifest", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		paths, err := filterFilesByTechnology(testLog, fs, sourceDir, []string{"java"})
		require.Error(t, err)
		assert.Nil(t, paths)
	})
}
