package lock

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/log"
	"github.com/go-logr/logr"
)

const (
	// DefaultStaleTimeout is the default duration after which a lock is considered stale.
	DefaultStaleTimeout = 5 * time.Minute

	filePerm600 fs.FileMode = 0o600
)

// FileLock represents a file-based lock with stale detection.
type FileLock struct {
	path         string
	staleTimeout time.Duration
	logger       logr.Logger
}

// New creates a new FileLock instance with the default stale timeout.
// The lock file will be created at the specified path when the lock is acquired.
func New(logger logr.Logger, filePath string) *FileLock {
	return &FileLock{
		path:         filePath,
		staleTimeout: DefaultStaleTimeout,
		logger:       logger,
	}
}

// WithStaleTimeout sets a custom stale timeout for the lock.
func (l *FileLock) WithStaleTimeout(timeout time.Duration) *FileLock {
	l.staleTimeout = timeout

	return l
}

// TryAcquire attempts to acquire the lock by using an exclusive file creation lock mechanism.
// This function does not guarantee an exclusive lock if a stale lock file is detected
// (e.g., if a process holding the lock crashed or was forcefully terminated),
// since removing the stale lock file and creating a new one is not atomic operation.
// The flock syscall is not used here because it is not supported on some NFS-mounted file systems.
//
// Therefore, this function provides a best-effort lock, but fake lock acquisition is possible!
// It is the caller's responsibility to handle this properly, for example, by deploying OneAgent binaries in
// a unique work directory per process to avoid write conflicts. Only one process should then atomically rename
// its work directory to the target location.
//
// Returns true if the lock was acquired, false if another process holds the lock.
// Returns an error if one occurred during lock acquisition.
func (l *FileLock) TryAcquire() (bool, error) {
	// If the lock file was not removed in a previous run and is now stale, remove it
	if l.isStale() {
		log.Debug(l.logger, "Detected stale lock file, removing it", "path", l.path)

		// As noted in the documentation: a race condition is still possible here because
		// checking for staleness and removing the lock file is not atomic operation.
		// If multiple processes detect the lock file as stale, one process might remove a new lock file created by
		// another process in the meantime.
		if err := os.Remove(l.path); err != nil && !os.IsNotExist(err) {
			l.logger.Info("Failed to remove stale lock file", "path", l.path, "error", err)
		}
	}

	// The os.O_CREATE|os.O_EXCL flags ensures the lock file is created atomically only if it doesn't exist.
	// The file's modification time is automatically set to the current time upon creation,
	// which is used for stale lock detection.
	fileLock, err := os.OpenFile(l.path, os.O_CREATE|os.O_EXCL, filePerm600)
	if err != nil {
		if os.IsExist(err) {
			l.logger.Info("Lock not acquired, lock file already exists", "path", l.path)

			return false, nil
		}

		return false, fmt.Errorf("failed to acquire lock: %w", err)
	}

	// Close immediately, only the file's existence and timestamp metadata matter - no write is needed.
	if err = fileLock.Close(); err != nil {
		l.logger.Error(err, "Failed to close lock file", "path", l.path)
	}

	log.Debug(l.logger, "Lock acquired successfully", "path", l.path)

	return true, nil
}

// Release removes the lock file to avoid the stale lock problem.
// It is the caller's responsibility to release the lock when no longer needed.
func (l *FileLock) Release() error {
	if err := os.Remove(l.path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to release lock: %w", err)
	}

	log.Debug(l.logger, "Lock released", "path", l.path)

	return nil
}

// isStale checks if the lock file exists and its modification time is older than staleTimeout.
func (l *FileLock) isStale() bool {
	fileInfo, err := os.Stat(l.path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug(l.logger, "Lock file does not exist", "path", l.path)

			return false
		}

		// File exists but can't be accessed, consider it stale to allow recovery
		l.logger.Info("Lock file exists but can't be accessed, considering it stale", "path", l.path, "error", err)

		return true
	}

	lockAge := time.Since(fileInfo.ModTime())
	if lockAge > l.staleTimeout {
		log.Debug(l.logger, "Lock file is stale", "age", lockAge.String(), "stale timeout", l.staleTimeout.String())

		return true
	}

	log.Debug(l.logger, "The lock file exists and is not stale", "path", l.path, "age", lockAge.String(), "stale timeout", l.staleTimeout.String())

	return false
}
