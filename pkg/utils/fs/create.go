package fs

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/spf13/afero"
)

func CreateFile(fs afero.Fs, path string, content string) error {
	err := fs.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}

	file, err := fs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
