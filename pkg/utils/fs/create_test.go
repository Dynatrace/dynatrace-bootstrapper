package fs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFile(t *testing.T) {
	t.Run("success, simple file", func(t *testing.T) {
		tmpDir := t.TempDir()

		expectedContent := "test\n\ntest"
		fileName := filepath.Join(tmpDir, "test.txt")

		err := CreateFile(fileName, expectedContent)
		require.NoError(t, err)

		content, err := os.ReadFile(fileName)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	})

	t.Run("success, nested file", func(t *testing.T) {
		tmpDir := t.TempDir()
		expectedContent := "test\n\ntest"

		fileName := filepath.Join(tmpDir, "folder", "inside", "test.txt")

		err := CreateFile(fileName, expectedContent)
		require.NoError(t, err)

		content, err := os.ReadFile(fileName)
		require.NoError(t, err)
		assert.Equal(t, expectedContent, string(content))
	})
}
