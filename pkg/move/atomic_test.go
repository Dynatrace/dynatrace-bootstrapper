package move

import (
	"errors"
	"path/filepath"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func mockCopyFuncWithAtomicCheck(t *testing.T, workFolder string, isSuccessful bool) copyFunc {
	t.Helper()

	return func(log logr.Logger, fs afero.Afero, _, target string) error {
		// according to the inner copyFunc, the target should be the workFolder
		// the actual target will be created outside the copyFunc by the atomic wrapper using fs.Rename
		require.Equal(t, workFolder, target)

		// the atomic wrapper should already have created the base workFolder
		exists, err := fs.DirExists(target)
		require.NoError(t, err)
		require.True(t, exists)

		if isSuccessful {
			file, err := fs.Create(filepath.Join(target, "test.txt"))
			require.NoError(t, err)

			_ = file.Close()

			return nil
		}

		return errors.New("some mock error")
	}
}
func TestAtomic(t *testing.T) {
	source := "/source"
	target := "/target"
	work := "/work"

	t.Run("success -> target is present", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		err := fs.MkdirAll(source, 0755)
		assert.NoError(t, err)

		atomicCopy := Atomic(work, mockCopyFuncWithAtomicCheck(t, work, true))

		err = atomicCopy(testLog, fs, source, target)
		assert.NoError(t, err)

		require.NotEqual(t, work, target)

		exists, err := fs.DirExists(work)
		assert.NoError(t, err)
		assert.False(t, exists)

		exists, err = fs.DirExists(target)
		assert.NoError(t, err)
		assert.True(t, exists)

		isEmpty, err := fs.IsEmpty(target)
		assert.NoError(t, err)
		assert.False(t, isEmpty)
	})
	t.Run("fail -> target is not present", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		atomicCopy := Atomic(work, mockCopyFuncWithAtomicCheck(t, work, false))

		err := atomicCopy(testLog, fs, source, target)
		assert.Error(t, err)
		assert.Equal(t, "some mock error", err.Error())

		require.NotEqual(t, work, target)

		exists, err := fs.DirExists(work)
		assert.NoError(t, err)
		assert.False(t, exists)

		exists, err = fs.DirExists(target)
		assert.NoError(t, err)
		assert.False(t, exists)
	})
}
