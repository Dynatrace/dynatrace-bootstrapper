package move

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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

func TestCopyFile(t *testing.T) {
	fs := afero.NewMemMapFs()
	source := "/source"
	target := "/target"

	err := fs.MkdirAll(source, 0755)
	require.NoError(t, err)

	err = afero.WriteFile(fs, filepath.Join(source, "file1.txt"), []byte("some content"), 0644)
	require.NoError(t, err)

	err = fs.MkdirAll(target, 0755)
	require.NoError(t, err)

	err = copyFile(fs, filepath.Join(source, "file1.txt"), filepath.Join(target, "file1.txt"))
	require.NoError(t, err)

	sourceContent, err := afero.ReadFile(fs, filepath.Join(source, "file1.txt"))
	require.NoError(t, err)
	assert.Equal(t, "some content", string(sourceContent))

	targetContent, err := afero.ReadFile(fs, filepath.Join(source, "file1.txt"))
	require.NoError(t, err)
	assert.Equal(t, "some content", string(targetContent))

	sourceFiles, err := afero.ReadDir(fs, source)
	require.NoError(t, err)

	targetFiles, err := afero.ReadDir(fs, target)
	require.NoError(t, err)
	require.Len(t, targetFiles, len(sourceFiles))
}
