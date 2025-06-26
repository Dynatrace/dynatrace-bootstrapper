package fs

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFile(t *testing.T) {
	t.Run("success, simple file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		expectedContent := "test\n\ntest"
		fileName := "test.txt"

		err := CreateFile(fs, fileName, expectedContent)
		require.NoError(t, err)

		content, err := afero.Afero{Fs: fs}.ReadFile(fileName)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	})

	t.Run("success, nested file", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		expectedContent := "test\n\ntest"
		fileName := "folder/inside/test.txt"

		err := CreateFile(fs, fileName, expectedContent)
		require.NoError(t, err)

		content, err := afero.Afero{Fs: fs}.ReadFile(fileName)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))

		stat, err := fs.Stat(filepath.Dir(fileName))
		require.NoError(t, err)
		assert.True(t, stat.IsDir())
	})
}
