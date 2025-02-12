package move

import (
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func mockCopyFunc(isSuccessful bool) copyFunc {
	return func(fs afero.Afero, _, targetFolder string) error {
		if isSuccessful {
			_ = fs.MkdirAll(targetFolder, 0755)
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

		atomicCopy := atomic(mockCopyFunc(true))

		err = atomicCopy(fs, sourceFolder, targetFolder)
		assert.NoError(t, err)

		exists, err := fs.DirExists(workFolder)
		assert.NoError(t, err)
		assert.False(t, exists)

		exists, err = fs.DirExists(targetFolder)
		assert.NoError(t, err)
		assert.True(t, exists)
	})
	t.Run("fail -> targetFolder is not present", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}
		sourceFolder = "/source"
		targetFolder = "/target"
		workFolder = "/work"

		atomicCopy := atomic(mockCopyFunc(false))

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
