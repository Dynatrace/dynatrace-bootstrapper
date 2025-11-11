package pmc

import (
	"os"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
)

func Create(log logr.Logger, srcPath, dstPath string, conf ruxit.ProcConf) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		log.Info("failed to open source file", "path", srcPath)

		return errors.WithStack(err)
	}

	defer func() { _ = srcFile.Close() }()

	srcInfo, err := srcFile.Stat()
	if err != nil {
		log.Info("failed to stat the source file", "path", srcPath)

		return err
	}

	srcConf, err := ruxit.FromConf(srcFile)
	if err != nil {
		log.Info("failed to parse source file to struct", "path", srcPath)

		return err
	}

	mergedConf := srcConf.Merge(conf)

	err = os.MkdirAll(filepath.Dir(dstPath), os.ModePerm)
	if err != nil {
		log.Info("failed to create destination dir", "path", filepath.Dir(filepath.Dir(dstPath)))

		return err
	}

	dstFile, err := os.OpenFile(dstPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		log.Info("failed to open destination file to write", "path", dstPath)

		return errors.WithStack(err)
	}

	defer func() { _ = dstFile.Close() }()

	_, err = dstFile.WriteString(mergedConf.ToString())
	if err != nil {
		log.Info("failed to write merged config into destination file", "path", dstPath)

		return errors.WithStack(err)
	}

	return nil
}
