package pmc

import (
	"os"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

const (
	RuxitAgentProcPath             = "agent/conf/ruxitagentproc.conf"
	originalCopyRuxitAgentProcPath = "agent/conf/_ruxitagentproc.conf"
)

func UpdateInPlace(log logr.Logger, fs afero.Fs, targetDir string, conf ruxit.ProcConf) error {
	log.Info("updating ruxitagentproc.conf", "targetDir", targetDir)
	destConfPath := filepath.Join(targetDir, RuxitAgentProcPath)
	sourceConfPath := filepath.Join(targetDir, originalCopyRuxitAgentProcPath)

	return safeMerge(log, fs, sourceConfPath, destConfPath, conf)
}

func safeMerge(log logr.Logger, fs afero.Fs, sourcePath, destPath string, conf ruxit.ProcConf) error {
	fileInfo, err := checkCopy(log, fs, sourcePath, destPath)
	if err != nil {
		log.Info("failed to create copy of original config", "path", destPath)

		return err
	}

	sourceFile, err := fs.Open(sourcePath)
	if err != nil {
		log.Info("failed to open source file", "path", sourcePath)

		return errors.WithStack(err)
	}

	defer func() { _ = sourceFile.Close() }()

	sourceConf, err := ruxit.FromConf(sourceFile)
	if err != nil {
		log.Info("failed to parse source file to struct", "path", sourcePath)

		return err
	}

	mergedConf := sourceConf.Merge(conf)

	destFile, err := fs.OpenFile(destPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, fileInfo.Mode())
	if err != nil {
		log.Info("failed to open destination file to write", "path", destPath)

		return errors.WithStack(err)
	}

	defer func() { _ = destFile.Close() }()

	_, err = destFile.Write([]byte(mergedConf.ToString()))
	if err != nil {
		log.Info("failed to write merged config into destination file", "path", destPath)

		return errors.WithStack(err)
	}

	return nil
}

func checkCopy(log logr.Logger, fs afero.Fs, sourcePath, destPath string) (os.FileInfo, error) {
	fileInfo, err := fs.Stat(sourcePath)
	if os.IsNotExist(err) {
		log.Info("saving copy of original config for transparency", "original", sourcePath, "copy", destPath)

		err := fs.Rename(destPath, sourcePath)
		if err != nil {
			return fileInfo, err
		}

		fileInfo, err = fs.Stat(sourcePath)
		if err != nil {
			return fileInfo, err
		}
	} else if err != nil {
		return fileInfo, err
	}

	return fileInfo, nil
}
