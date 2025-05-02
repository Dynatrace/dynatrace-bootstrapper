package pmc

import (
	"os"
	"path/filepath"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
	"github.com/go-logr/logr"
	"github.com/pkg/errors"
	"github.com/spf13/afero"
)


func Create(log logr.Logger, fs afero.Fs, source, destination string, conf ruxit.ProcConf) error {
	sourceFile, err := fs.Open(source)
	if err != nil {
		log.Info("failed to open source file", "path", source)

		return errors.WithStack(err)
	}

	defer func() { _ = sourceFile.Close() }()

	sourceInfo, err := sourceFile.Stat()
	if err != nil {
		log.Info("failed to stat the source file", "path", source)

		return err
	}

	sourceConf, err := ruxit.FromConf(sourceFile)
	if err != nil {
		log.Info("failed to parse source file to struct", "path", source)

		return err
	}

	mergedConf := sourceConf.Merge(conf)

	err = fs.MkdirAll(filepath.Dir(destination), os.ModePerm)
	if err != nil {
		log.Info("failed to create destination dir", "path", filepath.Dir(filepath.Dir(destination)))
		return err
	}

	destFile, err := fs.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, sourceInfo.Mode())
	if err != nil {
		log.Info("failed to open destination file to write", "path", destination)

		return errors.WithStack(err)
	}

	defer func() { _ = destFile.Close() }()

	_, err = destFile.Write([]byte(mergedConf.ToString()))
	if err != nil {
		log.Info("failed to write merged config into destination file", "path", destination)

		return errors.WithStack(err)
	}

	return nil
}
