package move

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockCopyFuncWithAtomicCheck(t *testing.T, isSuccessful bool) copyFunc {
	t.Helper()

	return func(fs afero.Afero, _, targetFolder string) error {
		// according to the inner copyFunc, the targetFolder should be the workFolder
		// the actual targetFolder will be created outside the copyFunc by the atomic wrapper using fs.Rename
		require.Equal(t, workFolder, targetFolder)

		// the atomic wrapper should already have created the base workFolder
		exists, err := fs.DirExists(targetFolder)
		require.NoError(t, err)
		require.True(t, exists)

		if isSuccessful {
			file, err := fs.Create(filepath.Join(targetFolder, "test.txt"))
			require.NoError(t, err)
			file.Close()

			return nil
		}

		return errors.New("some mock error")
	}
}
func TestAtomic(t *testing.T) {
	t.Run("success -> targetFolder is present", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		sourceFolder = "/source"
		targetFolder = "/target"
		workFolder = "/work"

		err := fs.MkdirAll(sourceFolder, 0755)
		assert.NoError(t, err)

		atomicCopy := atomic(mockCopyFuncWithAtomicCheck(t, true))

		err = atomicCopy(fs, sourceFolder, targetFolder)
		assert.NoError(t, err)

		isEmpty, err := fs.DirExists(workFolder)
		assert.NoError(t, err)
		assert.False(t, isEmpty)

		isEmpty, err = fs.DirExists(targetFolder)
		assert.NoError(t, err)
		assert.True(t, isEmpty)

		isEmpty, err = fs.IsEmpty(targetFolder)
		assert.NoError(t, err)
		assert.False(t, isEmpty)
	})
	t.Run("fail -> targetFolder is not present", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		sourceFolder = "/source"
		targetFolder = "/target"
		workFolder = "/work"

		atomicCopy := atomic(mockCopyFuncWithAtomicCheck(t, false))

		err := atomicCopy(fs, sourceFolder, targetFolder)
		assert.Error(t, err)
		assert.Equal(t, "some mock error", err.Error())

		exists, err := fs.DirExists(workFolder)
		assert.NoError(t, err)
		assert.False(t, exists)

		exists, err = fs.DirExists(targetFolder)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
