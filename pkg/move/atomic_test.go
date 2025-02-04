package move

import (
	"errors"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func mockCopyFunc(isSuccessful bool) copyFunc {
	return func(fs afero.Afero) error {
		if isSuccessful {
			fs.MkdirAll("/work/tmp", 0755)
			return nil
		}

		return errors.New("some mock error")
	}
}
func TestAtomic_NoError(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}
	sourceFolder = "/source"
	targetFolder = "/target"
	workFolder = "/work"

	fs.MkdirAll(sourceFolder, 0755)

	atomicCopy := atomic(mockCopyFunc(true))

	err := atomicCopy(fs)
	assert.NoError(t, err)

	exists, err := fs.DirExists(workFolder)
	assert.NoError(t, err)
	assert.False(t, exists)
}

func TestAtomic_Error(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}
	sourceFolder = "/source"
	targetFolder = "/target"
	workFolder = "/work"

	atomicCopy := atomic(mockCopyFunc(false))

	err := atomicCopy(fs)
	assert.Error(t, err)
	assert.Equal(t, "some mock error", err.Error())

	exists, err := fs.DirExists(workFolder)
	assert.NoError(t, err)
	assert.False(t, exists)
}
