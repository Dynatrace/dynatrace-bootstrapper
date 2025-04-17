package move

import (
	iofs "io/fs"
	"path/filepath"
	"regexp"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs/symlink"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const (
	// example match: 1.239.14.20220325-164521
	versionRegexp = `^(\d+)\.(\d+)\.(\d+)\.(\d+)-(\d+)$`
	binDir        = "/agent/bin"
	currentDir    = "current"
)

// CreateCurrentSymlink finds a the version of the CodeModule in the `targetDir` and creates a "current" symlink next to it.
// this is needed for the nginx use-case.
func CreateCurrentSymlink(log logr.Logger, fs afero.Fs, targetDir string) error {
	targetBinDir := filepath.Join(targetDir, binDir)

	relativeSymlinkPath, err := findVersionFromFS(log, fs, targetBinDir)
	if err != nil {
		log.Info("failed to get the version from the filesystem", "targetDir", targetDir)

		return err
	}

	return symlink.Create(log, fs, relativeSymlinkPath, filepath.Join(targetBinDir, currentDir))
}

func findVersionFromFS(log logr.Logger, fs afero.Fs, targetBinDir string) (string, error) {
	var version string

	aferoFs := afero.Afero{
		Fs: fs,
	}
	walkFiles := func(path string, info iofs.FileInfo, err error) error {
		if info == nil {
			log.Info(
				"version sub-dir does not exist in dir",
				"dir", targetBinDir)

			return iofs.ErrNotExist
		}

		if !info.IsDir() {
			return nil
		}

		folderName := filepath.Base(path)
		if regexp.MustCompile(versionRegexp).Match([]byte(folderName)) {
			log.Info("found version", "version", folderName)
			version = folderName

			return iofs.ErrExist
		}

		return nil
	}

	err := aferoFs.Walk(targetBinDir, walkFiles)
	if errors.Is(err, iofs.ErrNotExist) {
		return "", errors.WithStack(err)
	}

	return version, nil
}
