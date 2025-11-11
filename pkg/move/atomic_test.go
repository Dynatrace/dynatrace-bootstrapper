package move

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func mockCopyFuncWithAtomicCheck(t *testing.T, workFolder string, isSuccessful bool) CopyFunc {
	t.Helper()

	return func(_ logr.Logger, _, target string) error {
		// according to the inner copyFunc, the target should be the workFolder
		// the actual target will be created outside the copyFunc by the atomic wrapper using fs.Rename
		require.Equal(t, workFolder, target)

		// the atomic wrapper should already have created the base workFolder
		assert.DirExists(t, target)

		if isSuccessful {
			file, err := os.Create(filepath.Join(target, "test.txt"))
			require.NoError(t, err)

			_ = file.Close()

			return nil
		}

		return errors.New("some mock error")
	}
}
func TestAtomic(t *testing.T) {
	t.Run("success -> target is present", func(t *testing.T) {
		tmpDir := t.TempDir()
		source := filepath.Join(tmpDir, "source")
		target := filepath.Join(tmpDir, "target")
		work := filepath.Join(tmpDir, "work")

		err := os.MkdirAll(source, 0755)
		require.NoError(t, err)

		atomicCopy := Atomic(work, mockCopyFuncWithAtomicCheck(t, work, true))

		err = atomicCopy(testLog, source, target)
		require.NoError(t, err)

		require.NotEqual(t, work, target)

		assert.NoDirExists(t, work)
		require.DirExists(t, target)

		files, err := os.ReadDir(target)
		require.NoError(t, err)
		require.NotEmpty(t, files)
	})
	t.Run("fail -> target is not present", func(t *testing.T) {
		tmpDir := t.TempDir()
		source := filepath.Join(tmpDir, "source")
		target := filepath.Join(tmpDir, "target")
		work := filepath.Join(tmpDir, "work")

		atomicCopy := Atomic(work, mockCopyFuncWithAtomicCheck(t, work, false))

		err := atomicCopy(testLog, source, target)
		require.Error(t, err)
		assert.Equal(t, "some mock error", err.Error())

		require.NotEqual(t, work, target)

		assert.NoDirExists(t, work)
		assert.NoDirExists(t, target)
	})
}
