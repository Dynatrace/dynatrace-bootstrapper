package move

import (
	"fmt"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute(t *testing.T) {
	t.Run("package global vars are used", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		// Create source directory and files
		sourceDir := "/source"
		targetDir := "/target"
		workDir := "/work"
		_ = fs.MkdirAll(sourceDir, 0755)
		_ = afero.WriteFile(fs, sourceDir+"/file1.txt", []byte("file1 content"), 0644)
		_ = afero.WriteFile(fs, sourceDir+"/file2.txt", []byte("file2 content"), 0644)

		sourceFolder = sourceDir
		targetFolder = targetDir
		workFolder = workDir

		err := Execute(fs)
		require.NoError(t, err)

		// Check if the target directory and files exist
		exists, err := afero.DirExists(fs, targetDir)
		require.NoError(t, err)
		assert.True(t, exists)

		file1Exists, err := afero.Exists(fs, targetDir+"/file1.txt")
		assert.NoError(t, err)
		assert.True(t, file1Exists)

		file2Exists, err := afero.Exists(fs, targetDir+"/file2.txt")
		assert.NoError(t, err)
		assert.True(t, file2Exists)

		// Check the content of the copied files
		content, err := afero.ReadFile(fs, targetDir+"/file1.txt")
		assert.NoError(t, err)
		assert.Equal(t, "file1 content", string(content))

		content, err = afero.ReadFile(fs, targetDir+"/file2.txt")
		assert.NoError(t, err)
		assert.Equal(t, "file2 content", string(content))

		// Check the cleanup happened of the copied files
		exists, err = afero.DirExists(fs, workDir)
		require.NoError(t, err)
		assert.False(t, exists)
	})
}

func TestCopyFolder(t *testing.T) {
	fs := afero.NewMemMapFs()
	src := "/src"
	err := fs.MkdirAll(src, 0755)
	require.NoError(t, err)

	err = afero.WriteFile(fs, filepath.Join(src, "file1.txt"), []byte("Hello"), 0644)
	require.NoError(t, err)

	err = fs.MkdirAll(filepath.Join(src, "subdir"), 0755)
	require.NoError(t, err)

	err = afero.WriteFile(fs, filepath.Join(src, "subdir", "file2.txt"), []byte("World"), 0644)
	require.NoError(t, err)

	dst := "/dst"
	err = fs.MkdirAll(dst, 0755)
	require.NoError(t, err)

	err = copyFolder(fs, src, dst)
	require.NoError(t, err)

	srcFiles, err := afero.ReadDir(fs, src)
	require.NoError(t, err)
	dstFiles, err := afero.ReadDir(fs, dst)
	require.NoError(t, err)
	require.Len(t, dstFiles, len(srcFiles))

	checkFolder(t, fs, src, dst)
}

func TestCopyFolderWithTechnologyFiltering(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	sourceDir := "/source"
	targetDir := "/target"

	sourceFolder = sourceDir
	targetFolder = targetDir

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
		technology = "java"
		err := copyFolder(fs, sourceDir, targetDir)
		require.NoError(t, err)

		assertFileExists(t, fs, filepath.Join(targetDir, "fileA1.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileA2.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileB1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileC1.txt"))
	})
	t.Run("copy with multiple technology filter", func(t *testing.T) {
		technology = "java,python"
		err := copyFolder(fs, sourceDir, targetDir)
		require.NoError(t, err)

		assertFileExists(t, fs, filepath.Join(targetDir, "fileA1.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileA2.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileB1.txt"))
		assertFileNotExists(t, fs, filepath.Join(targetDir, "fileC1.txt"))
	})
	t.Run("copy with invalid technology filter", func(t *testing.T) {
		technology = "php"
		err := copyFolder(fs, sourceDir, targetDir)
		require.NoError(t, err)

		assertFileExists(t, fs, filepath.Join(targetDir, "fileA1.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileA2.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileB1.txt"))
		assertFileExists(t, fs, filepath.Join(targetDir, "fileC1.txt"))
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
		paths, err := filterFilesByTechnology(fs, sourceDir, []string{"java"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			filepath.Join(sourceDir, "fileA1.txt"),
			filepath.Join(sourceDir, "fileA2.txt"),
		}, paths)
	})
	t.Run("filter multiple technologies", func(t *testing.T) {
		paths, err := filterFilesByTechnology(fs, sourceDir, []string{"java", "python"})
		require.NoError(t, err)
		assert.ElementsMatch(t, []string{
			filepath.Join(sourceDir, "fileA1.txt"),
			filepath.Join(sourceDir, "fileA2.txt"),
			filepath.Join(sourceDir, "fileB1.txt"),
		}, paths)
	})
	t.Run("not filter non-existing technology", func(t *testing.T) {
		paths, err := filterFilesByTechnology(fs, sourceDir, []string{"php"})
		require.NoError(t, err)
		assert.Empty(t, paths)
	})
	t.Run("filter with missing manifest", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		paths, err := filterFilesByTechnology(fs, sourceDir, []string{"java"})
		require.Error(t, err)
		assert.Nil(t, paths)
	})
}

func checkFolder(t *testing.T, fs afero.Fs, src, dst string) {
	srcFiles, err := afero.ReadDir(fs, src)
	require.NoError(t, err)
	dstFiles, err := afero.ReadDir(fs, dst)
	require.NoError(t, err)
	require.Len(t, dstFiles, len(srcFiles))

	for i := range srcFiles {
		srcName := srcFiles[i].Name()
		dstName := dstFiles[i].Name()
		require.Equal(t, srcName, dstName)

		srcPath := filepath.Join(src, srcName)
		dstPath := filepath.Join(dst, dstName)

		srcInfo, err := fs.Stat(srcPath)
		require.NoError(t, err)

		dstInfo, err := fs.Stat(dstPath)
		require.NoError(t, err)

		assert.Equal(t, srcInfo.Mode(), dstInfo.Mode())

		if srcInfo.IsDir() {
			assert.True(t, dstInfo.IsDir())
			checkFolder(t, fs, srcPath, dstPath)
		} else {
			srcData, err := afero.ReadFile(fs, srcPath)
			require.NoError(t, err)

			dstData, err := afero.ReadFile(fs, dstPath)
			require.NoError(t, err)

			assert.Equal(t, srcData, dstData)
		}
	}
}
