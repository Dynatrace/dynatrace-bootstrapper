package move

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testSourceDir = "/source"

func TestCopyFolderWithTechnologyFiltering(t *testing.T) {
	tmpDir := t.TempDir()

	sourceDir := filepath.Join(tmpDir, testSourceDir)
	targetDir := filepath.Join(tmpDir, "/target")

	_ = os.MkdirAll(sourceDir, 0755)
	_ = os.MkdirAll(targetDir, 0755)

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

	_ = os.WriteFile(filepath.Join(sourceDir, "manifest.json"), []byte(manifestContent), 0600)
	_ = os.WriteFile(filepath.Join(sourceDir, "fileA1.txt"), []byte("java a1"), 0600)
	_ = os.WriteFile(filepath.Join(sourceDir, "fileA2.txt"), []byte("java a2"), 0600)
	_ = os.WriteFile(filepath.Join(sourceDir, "fileB1.txt"), []byte("python b1"), 0600)
	_ = os.WriteFile(filepath.Join(sourceDir, "fileC1.txt"), []byte("unrelated"), 0600)

	t.Run("copy with single technology filter", func(t *testing.T) {
		t.Cleanup(func() {
			_ = os.RemoveAll(targetDir)
			_ = os.MkdirAll(targetDir, 0755)
		})

		technology := "java"
		err := CopyByTechnology(testLog, sourceDir, targetDir, technology)
		require.NoError(t, err)

		assert.FileExists(t, filepath.Join(targetDir, "fileA1.txt"))
		assert.FileExists(t, filepath.Join(targetDir, "fileA2.txt"))
		assert.NoFileExists(t, filepath.Join(targetDir, "fileB1.txt"))
		assert.NoFileExists(t, filepath.Join(targetDir, "fileC1.txt"))
	})
	t.Run("copy with multiple technology filter", func(t *testing.T) {
		t.Cleanup(func() {
			_ = os.RemoveAll(targetDir)
			_ = os.MkdirAll(targetDir, 0755)
		})

		technology := "java,python"
		err := CopyByTechnology(testLog, sourceDir, targetDir, technology)
		require.NoError(t, err)

		assert.FileExists(t, filepath.Join(targetDir, "fileA1.txt"))
		assert.FileExists(t, filepath.Join(targetDir, "fileA2.txt"))
		assert.FileExists(t, filepath.Join(targetDir, "fileB1.txt"))
		assert.NoFileExists(t, filepath.Join(targetDir, "fileC1.txt"))
	})
	t.Run("copy with multiple technology filter with whitespace", func(t *testing.T) {
		t.Cleanup(func() {
			_ = os.RemoveAll(targetDir)
			_ = os.MkdirAll(targetDir, 0755)
		})

		technology := "java, python"
		err := CopyByTechnology(testLog, sourceDir, targetDir, technology)
		require.NoError(t, err)

		assert.FileExists(t, filepath.Join(targetDir, "fileA1.txt"))
		assert.FileExists(t, filepath.Join(targetDir, "fileA2.txt"))
		assert.FileExists(t, filepath.Join(targetDir, "fileB1.txt"))
		assert.NoFileExists(t, filepath.Join(targetDir, "fileC1.txt"))
	})
	t.Run("copy with invalid technology filter", func(t *testing.T) {
		t.Cleanup(func() {
			_ = os.RemoveAll(targetDir)
			_ = os.MkdirAll(targetDir, 0755)
		})

		technology := "php"
		err := CopyByTechnology(testLog, sourceDir, targetDir, technology)
		require.NoError(t, err)

		assert.NoFileExists(t, filepath.Join(targetDir, "fileA1.txt"))
		assert.NoFileExists(t, filepath.Join(targetDir, "fileA2.txt"))
		assert.NoFileExists(t, filepath.Join(targetDir, "fileB1.txt"))
		assert.NoFileExists(t, filepath.Join(targetDir, "fileC1.txt"))
	})
}

func TestCopyByList(t *testing.T) {
	tmpDir := t.TempDir()

	sourceDir := filepath.Join(tmpDir, "source")
	require.NoError(t, os.Mkdir(sourceDir, os.ModePerm))

	dirs := []string{
		"folder",
		filepath.Join("folder", "sub"),
		filepath.Join("folder", "sub", "child"),
	}

	dirModes := []os.FileMode{
		0744, // rwxr-xr-x
		0740, // rwxr--r--
		0700, // rwx------
	}

	filesNames := []string{
		"f1.txt",
		"runtime",
		"log",
	}

	fileModes := []os.FileMode{
		0744, // rwxr-xr-x
		0740, // rwxr--r--
		0700, // rwx------
	}

	// create an FS where there are multiple sub dirs and files, each with their own file modes
	for i := range dirs {
		dir := filepath.Join(sourceDir, dirs[i])

		err := os.Mkdir(dir, dirModes[i])
		require.NoError(t, err)
		err = os.WriteFile(filepath.Join(dir, filesNames[i]), fmt.Appendf(nil, "%d", i), fileModes[i])
		require.NoError(t, err)
	}

	// reverse the list, so the longest path is the first
	fileList := []string{}
	for i := len(dirs) - 1; i >= 0; i-- {
		fileList = append(fileList, filepath.Join(dirs[i], filesNames[i]))
	}

	targetDir := filepath.Join(tmpDir, "target")

	err := copyByList(testLog, sourceDir, targetDir, fileList)
	require.NoError(t, err)

	for i := range dirs {
		targetDirStat, err := os.Stat(filepath.Join(targetDir, dirs[i]))
		require.NoError(t, err)
		assert.Equal(t, dirModes[i].Perm().String(), targetDirStat.Mode().Perm().String(), targetDirStat.Name())

		sourceDirStat, err := os.Stat(filepath.Join(sourceDir, dirs[i]))
		require.NoError(t, err)
		assert.Equal(t, targetDirStat.Mode().String(), sourceDirStat.Mode().String(), sourceDirStat.Name())

		targetFileStat, err := os.Stat(filepath.Join(targetDir, dirs[i], filesNames[i]))
		require.NoError(t, err)
		assert.Equal(t, fileModes[i].String(), targetFileStat.Mode().String(), targetFileStat.Name())

		sourceFileStat, err := os.Stat(filepath.Join(sourceDir, dirs[i], filesNames[i]))
		require.NoError(t, err)
		assert.Equal(t, targetFileStat.Mode().String(), sourceFileStat.Mode().String(), sourceFileStat.Name())
	}
}

func TestFilterFilesByTechnology(t *testing.T) {
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

	t.Run("filter single technology", func(t *testing.T) {
		tmpDir := t.TempDir()

		sourceDir := filepath.Join(tmpDir, testSourceDir)
		_ = os.MkdirAll(sourceDir, 0755)
		_ = os.WriteFile(filepath.Join(sourceDir, "manifest.json"), []byte(manifestContent), 0600)
		paths, err := filterFilesByTechnology(testLog, sourceDir, []string{"java"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			"fileA1.txt",
			"fileA2.txt",
		}, paths)
	})
	t.Run("filter multiple technologies", func(t *testing.T) {
		tmpDir := t.TempDir()

		sourceDir := filepath.Join(tmpDir, testSourceDir)
		_ = os.MkdirAll(sourceDir, 0755)
		_ = os.WriteFile(filepath.Join(sourceDir, "manifest.json"), []byte(manifestContent), 0600)
		paths, err := filterFilesByTechnology(testLog, sourceDir, []string{"java", "python"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			"fileA1.txt",
			"fileA2.txt",
			"fileB1.txt",
		}, paths)
	})
	t.Run("filter multiple technologies with white spaces", func(t *testing.T) {
		tmpDir := t.TempDir()

		sourceDir := filepath.Join(tmpDir, testSourceDir)
		_ = os.MkdirAll(sourceDir, 0755)
		_ = os.WriteFile(filepath.Join(sourceDir, "manifest.json"), []byte(manifestContent), 0600)
		paths, err := filterFilesByTechnology(testLog, sourceDir, []string{"java ", " python "})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			"fileA1.txt",
			"fileA2.txt",
			"fileB1.txt",
		}, paths)
	})
	t.Run("not filter non-existing technology", func(t *testing.T) {
		tmpDir := t.TempDir()

		sourceDir := filepath.Join(tmpDir, testSourceDir)
		_ = os.MkdirAll(sourceDir, 0755)
		_ = os.WriteFile(filepath.Join(sourceDir, "manifest.json"), []byte(manifestContent), 0600)

		paths, err := filterFilesByTechnology(testLog, sourceDir, []string{"php"})
		require.NoError(t, err)
		assert.Empty(t, paths)
	})
	t.Run("filter with missing manifest", func(t *testing.T) {
		tmpDir := t.TempDir()

		sourceDir := filepath.Join(tmpDir, testSourceDir)
		_ = os.MkdirAll(sourceDir, 0755)

		paths, err := filterFilesByTechnology(testLog, sourceDir, []string{"java"})
		require.Error(t, err)
		assert.Nil(t, paths)
	})
}
