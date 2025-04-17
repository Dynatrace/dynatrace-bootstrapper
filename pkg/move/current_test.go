package move

import (
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestCreateCurrentSymlink(t *testing.T) {
	testPath := "/test"
	expectedVersion := "1.239.14.20220325-164521"

	t.Run("no fail if version file exists", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		_ = fs.WriteFile(filepath.Join(testPath, InstallerVersionFilePath), []byte(expectedVersion), 0644)

		err := CreateCurrentSymlink(testLog, fs, testPath)
		require.NoError(t, err)
	})
	t.Run("fail if version file is missing", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		err := CreateCurrentSymlink(testLog, fs, testPath)
		require.Error(t, err)
	})
}
