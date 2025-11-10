package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestCopyFolder(t *testing.T) {
	tmpDir := t.TempDir()

	src := filepath.Join(tmpDir, "src")
	err := os.MkdirAll(src, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(src, "file1.txt"), []byte("Hello"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(filepath.Join(src, "subdir"), 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(src, "subdir", "file2.txt"), []byte("World"), 0644)
	require.NoError(t, err)

	dst := filepath.Join(tmpDir, "dst")
	err = os.MkdirAll(dst, 0755)
	require.NoError(t, err)

	err = CopyFolder(testLog, src, dst)
	require.NoError(t, err)

	srcFiles, err := os.ReadDir(src)
	require.NoError(t, err)
	dstFiles, err := os.ReadDir(dst)
	require.NoError(t, err)
	require.Len(t, dstFiles, len(srcFiles))

	checkFolder(t, src, dst)
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	source := filepath.Join(tmpDir, "/source")
	target := filepath.Join(tmpDir, "/target")

	err := os.MkdirAll(source, 0755)
	require.NoError(t, err)

	err = os.WriteFile(filepath.Join(source, "file1.txt"), []byte("some content"), 0644)
	require.NoError(t, err)

	err = os.MkdirAll(target, 0755)
	require.NoError(t, err)

	err = CopyFile(filepath.Join(source, "file1.txt"), filepath.Join(target, "file1.txt"))
	require.NoError(t, err)

	sourceContent, err := os.ReadFile(filepath.Join(source, "file1.txt"))
	require.NoError(t, err)
	assert.Equal(t, "some content", string(sourceContent))

	targetContent, err := os.ReadFile(filepath.Join(source, "file1.txt"))
	require.NoError(t, err)
	assert.Equal(t, "some content", string(targetContent))

	sourceFiles, err := os.ReadDir(source)
	require.NoError(t, err)

	targetFiles, err := os.ReadDir(target)
	require.NoError(t, err)
	require.Len(t, targetFiles, len(sourceFiles))
}

func checkFolder(t *testing.T, src, dst string) {
	t.Helper()

	srcFiles, err := os.ReadDir(src)
	require.NoError(t, err)
	dstFiles, err := os.ReadDir(dst)
	require.NoError(t, err)
	require.Len(t, dstFiles, len(srcFiles))

	for i := range srcFiles {
		srcName := srcFiles[i].Name()
		dstName := dstFiles[i].Name()
		require.Equal(t, srcName, dstName)

		srcPath := filepath.Join(src, srcName)
		dstPath := filepath.Join(dst, dstName)

		srcInfo, err := os.Stat(srcPath)
		require.NoError(t, err)

		dstInfo, err := os.Stat(dstPath)
		require.NoError(t, err)

		assert.Equal(t, srcInfo.Mode(), dstInfo.Mode())

		if srcInfo.IsDir() {
			assert.True(t, dstInfo.IsDir())
			checkFolder(t, srcPath, dstPath)
		} else {
			srcData, err := os.ReadFile(srcPath)
			require.NoError(t, err)

			dstData, err := os.ReadFile(dstPath)
			require.NoError(t, err)

			assert.Equal(t, srcData, dstData)
		}
	}
}
