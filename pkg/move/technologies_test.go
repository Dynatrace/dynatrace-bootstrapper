package move

import (
	"fmt"
	"io/fs"
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

		technology := "java"
		err := CopyByTechnology(testLog, fs, sourceDir, targetDir, technology)
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

		technology := "java,python"
		err := CopyByTechnology(testLog, fs, sourceDir, targetDir, technology)
		require.NoError(t, err)

		assertFileExists(t, fs, filepath.Join(targetDir, "fileA1.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileA2.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileB1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileC1.txt"))
	})
	t.Run("copy with multiple technology filter with whitespace", func(t *testing.T) {
		t.Cleanup(func() {
			_ = fs.RemoveAll(targetDir)
			_ = fs.MkdirAll(targetDir, 0755)
		})

		technology := "java, python"
		err := CopyByTechnology(testLog, fs, sourceDir, targetDir, technology)
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

		technology := "php"
		err := CopyByTechnology(testLog, fs, sourceDir, targetDir, technology)
		require.NoError(t, err)

		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileA1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileA2.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileB1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileC1.txt"))
	})
}

func TestCopyByList(t *testing.T) {
	dirs := []string{
		"./folder",
		"./folder/sub",
		"./folder/sub/child",
	}
	dirModes := []fs.FileMode{
		0777,
		0776,
		0775,
	}

	filesNames := []string{
		"f1.txt",
		"runtime",
		"log",
	}
	fileModes := []fs.FileMode{
		0764,
		0773,
		0772,
	}

	fs := afero.Afero{Fs: afero.NewMemMapFs()}
	// create an FS where there are multiple sub dirs and files, each with their own file modes
	for i := range len(dirs) {
		err := fs.Mkdir(dirs[i], dirModes[i])
		require.NoError(t, err)
		err = fs.WriteFile(filepath.Join(dirs[i], filesNames[i]), []byte(fmt.Sprintf("%d", i)), fileModes[i])
		require.NoError(t, err)
	}

	// reverse the list, so the longest path is the first
	fileList := []string{}
	for i := len(dirs) - 1; i >= 0; i-- {
		fileList = append(fileList, filepath.Join(dirs[i], filesNames[i]))
	}

	targetDir := "./target"

	err := copyByList(testLog, fs, "./", targetDir, fileList)
	require.NoError(t, err)

	for i := range len(dirs) {
		targetStat, err := fs.Stat(filepath.Join(targetDir, dirs[i]))
		require.NoError(t, err)
		assert.Equal(t, dirModes[i], targetStat.Mode().Perm(), targetStat.Name())

		sourceStat, err := fs.Stat(dirs[i])
		require.NoError(t, err)
		assert.Equal(t, sourceStat.Mode(), targetStat.Mode(), targetStat.Name())

		targetStat, err = fs.Stat(filepath.Join(targetDir, dirs[i], filesNames[i]))
		require.NoError(t, err)
		assert.Equal(t, fileModes[i], targetStat.Mode().Perm(), targetStat.Name())

		sourceStat, err = fs.Stat(filepath.Join(dirs[i], filesNames[i]))
		require.NoError(t, err)
		assert.Equal(t, sourceStat.Mode(), targetStat.Mode(), targetStat.Name())
	}
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

	t.Run("filter single technology", func(t *testing.T) {
		paths, err := filterFilesByTechnology(testLog, fs, sourceDir, []string{"java"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			filepath.Join("fileA1.txt"),
			filepath.Join("fileA2.txt"),
		}, paths)
	})
	t.Run("filter multiple technologies", func(t *testing.T) {
		paths, err := filterFilesByTechnology(testLog, fs, sourceDir, []string{"java", "python"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			filepath.Join("fileA1.txt"),
			filepath.Join("fileA2.txt"),
			filepath.Join("fileB1.txt"),
		}, paths)
	})
	t.Run("filter multiple technologies with white spaces", func(t *testing.T) {
		paths, err := filterFilesByTechnology(testLog, fs, sourceDir, []string{"java ", " python "})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			filepath.Join("fileA1.txt"),
			filepath.Join("fileA2.txt"),
			filepath.Join("fileB1.txt"),
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

func assertFileExists(t *testing.T, fs afero.Fs, path string) {
	t.Helper()

	exists, err := afero.Exists(fs, path)
	assert.NoError(t, err)
	assert.True(t, exists, fmt.Sprintf("file should exist: %s", path))
}
func assertFileNotExists(t *testing.T, fs afero.Fs, path string) {
	t.Helper()

	exists, err := afero.Exists(fs, path)
	assert.NoError(t, err)
	assert.False(t, exists, fmt.Sprintf("file should not exist: %s", path))
}
