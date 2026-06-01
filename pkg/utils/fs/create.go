package fs

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// MostlyReadonlyFilePerm is only readonly for everyone but the owner.
// The owner needs to be able to write it, so we can seamlessly handle node restarts, where the files are not cleaned up.
const MostlyReadonlyFilePerm os.FileMode = 0644

func CreateFile(path string, content string) error {
	return createFileImpl(path, content, os.ModePerm)
}

func CreateReadOnlyFile(path string, content string) error {
	return createFileImpl(path, content, MostlyReadonlyFilePerm)
}

func createFileImpl(path string, content string, mode os.FileMode) error {
	// all created folders need to be writable, as the agent may write into them
	err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
	if err != nil {
		return errors.WithStack(err)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, mode)
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
