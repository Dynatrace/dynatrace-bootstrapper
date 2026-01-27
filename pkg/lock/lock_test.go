package lock

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/tests"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const lockFile = "test.lock"

func TestExecute(t *testing.T) {
	t.Run("Create new FileLock instance", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		fileLock := New(logger, lockFilePath)

		assert.Equal(t, lockFilePath, fileLock.path)
		assert.Equal(t, DefaultStaleTimeout, fileLock.staleTimeout)
		assert.Equal(t, logger, fileLock.logger)
	})

	t.Run("Set custom stale timeout", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		customStaleTimeout := 5 * time.Minute

		fileLock := New(logger, lockFilePath).WithStaleTimeout(customStaleTimeout)
		assert.Equal(t, customStaleTimeout, fileLock.staleTimeout)
	})

	t.Run("TryAcquire acquires lock successfully", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		fileLock := New(logger, lockFilePath)

		acquired, err := fileLock.TryAcquire()
		require.NoError(t, err)
		assert.True(t, acquired)
		// cleanup
		defer func() {
			require.NoError(t, fileLock.Release())
		}()

		// verify the lock file exists
		info, err := os.Stat(lockFilePath)
		require.NoError(t, err)
		assert.NotNil(t, info)
		assert.False(t, info.IsDir())
	})

	t.Run("TryAcquire returns `false` when lock already held", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		fileLock := New(logger, lockFilePath)

		// the first acquisition should succeed
		acquired, err := fileLock.TryAcquire()
		require.NoError(t, err)
		require.True(t, acquired)
		// cleanup
		defer func() {
			require.NoError(t, fileLock.Release())
		}()

		// the second acquisition should fail
		acquired, err = fileLock.TryAcquire()
		require.NoError(t, err)
		require.False(t, acquired)
	})

	t.Run("TryAcquire removes stale lock and acquires the new lock", func(t *testing.T) {
		logger, logsObserver := tests.NewTestLogger()

		// create a lock file
		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		staleLock := New(logger, lockFilePath)
		acquired, err := staleLock.TryAcquire()
		require.NoError(t, err)
		assert.True(t, acquired)

		// wait at least `customStaleTimeout` to simulate a stale lock scenario
		customStaleTimeout := 500 * time.Millisecond
		time.Sleep(customStaleTimeout + 100*time.Millisecond)

		// try to acquire lock again, should remove the stale lock and acquire the new lock
		fileLock := New(logger, lockFilePath).WithStaleTimeout(customStaleTimeout)
		acquired, err = fileLock.TryAcquire()
		require.NoError(t, err)
		assert.True(t, acquired)
		// cleanup
		defer func() {
			require.NoError(t, fileLock.Release())
		}()

		tests.RequireLogMessage(t, logsObserver, "Detected stale lock file, removing it", "path", lockFilePath)
	})

	t.Run("TryAcquire does not remove the fresh lock", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		fileLock1 := New(logger, lockFilePath)
		acquired, err := fileLock1.TryAcquire()
		require.NoError(t, err)
		assert.True(t, acquired)
		// cleanup
		defer func() {
			require.NoError(t, fileLock1.Release())
		}()

		// Create a new FileLock instance and try to acquire the lock again.
		// It should fail since the lock file's timestamp hasn't reached the stale timeout yet.
		staleTimeout := 5 * time.Minute
		fileLock2 := New(logger, lockFilePath).WithStaleTimeout(staleTimeout)
		acquired, err = fileLock2.TryAcquire()
		require.NoError(t, err)
		assert.False(t, acquired)
	})

	t.Run("TryAcquire fails when lock directory does not exist", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := "/nonexisting/directory/test.lock"
		fileLock := New(logger, lockFilePath)

		acquired, err := fileLock.TryAcquire()
		require.ErrorIs(t, err, syscall.ENOENT)
		require.ErrorContains(t, err, "failed to acquire lock: open "+lockFilePath+": no such file or directory")
		assert.False(t, acquired)
	})

	t.Run("Release successfully removes the lock file", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		fileLock := New(logger, lockFilePath)

		acquired, err := fileLock.TryAcquire()
		require.NoError(t, err)
		require.True(t, acquired)

		// cleanup
		err = fileLock.Release()
		require.NoError(t, err)

		// verify the lock file no longer exists
		_, err = os.Stat(lockFilePath)
		assert.True(t, os.IsNotExist(err))
	})

	t.Run("Release returns no error even when releasing non-acquired lock", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		fileLock := New(logger, lockFilePath)

		// cleanup
		defer func() {
			require.NoError(t, fileLock.Release())
		}()
	})

	t.Run("IsStale returns `false` when file does not exist", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		fileLock := New(logger, lockFilePath)

		assert.False(t, fileLock.isStale())
	})

	t.Run("IsStale returns `false` for fresh lock", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		fileLock := New(logger, lockFilePath)

		acquired, err := fileLock.TryAcquire()
		require.NoError(t, err)
		assert.True(t, acquired)
		// cleanup
		defer func() {
			require.NoError(t, fileLock.Release())
		}()

		// the default stale timeout is 10 minutes, so the acquire lock should be fresh
		assert.False(t, fileLock.isStale())
	})

	t.Run("IsStale returns `true` for stale lock", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)
		customStaleTimeout := 500 * time.Millisecond

		fileLock := New(logger, lockFilePath).WithStaleTimeout(customStaleTimeout)
		acquired, err := fileLock.TryAcquire()
		require.NoError(t, err)
		assert.True(t, acquired)
		// cleanup
		defer func() {
			require.NoError(t, fileLock.Release())
		}()

		// wait at least `customStaleTimeout` to simulate a stale lock scenario
		time.Sleep(customStaleTimeout + 100*time.Millisecond)
		assert.True(t, fileLock.isStale())
	})

	t.Run("In a concurrent environment, only one goroutine can acquire the lock.", func(t *testing.T) {
		logger, _ := tests.NewTestLogger()

		lockFilePath := filepath.Join(t.TempDir(), lockFile)

		const numGoroutines = 150

		results := make(chan bool, numGoroutines)

		for range numGoroutines {
			go func() {
				fileLock := New(logger, lockFilePath)

				acquired, err := fileLock.TryAcquire()
				if err != nil {
					results <- false

					return
				}

				results <- acquired
			}()
		}

		acquiredCount := 0

		for range numGoroutines {
			if <-results {
				acquiredCount++
			}
		}
		// cleanup
		require.NoError(t, os.Remove(lockFilePath))

		// only one goroutine should have acquired the lock
		assert.Equal(t, 1, acquiredCount)
	})
}
