package pmc

import (
	"os"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/oneagent/pmc/ruxit"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/fs"
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

	srcConf, err := ruxit.FromConf(srcFile)
	if err != nil {
		log.Info("failed to parse source file to struct", "path", srcPath)

		return err
	}

	mergedConf := srcConf.Merge(conf)

	return fs.CreateReadOnlyFile(dstPath, mergedConf.ToString())
}
