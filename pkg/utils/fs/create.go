package fs

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

func CreateFile(path string, content string) error {
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}
