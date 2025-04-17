package cmd

import (
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestBootstrapper(t *testing.T) {
	t.Run("should validate required flags - missing flags -> error", func(t *testing.T) {
		cmd := New(afero.NewMemMapFs())

		err := cmd.Execute()

		require.Error(t, err)
	})
	t.Run("should validate required flags - present flags -> no error", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_ = afero.WriteFile(fs, filepath.Join("./", move.InstallerVersionFilePath), []byte("123"), 0644)

		cmd := New(fs)
		cmd.SetArgs([]string{"--source", "./", "--target", "./"})

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("--suppress-error=true -> no error", func(t *testing.T) {
		cmd := New(afero.NewMemMapFs())
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\"})

		err := cmd.Execute()

		require.Error(t, err)

		cmd = New(afero.NewMemMapFs())
		// Note: we can't skip/mask the validation of the required flags
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\", "--suppress-error"})

		err = cmd.Execute()

		require.NoError(t, err)

		cmd = New(afero.NewMemMapFs())
		// Note: we can't skip/mask the validation of the required flags
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\", "--suppress-error", "true"})

		err = cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("should allow unknown flags -> no error", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		_ = afero.WriteFile(fs, filepath.Join("./", move.InstallerVersionFilePath), []byte("123"), 0644)

		cmd := New(fs)
		cmd.SetArgs([]string{"--source", "./", "--target", "./", "--unknown", "--flag", "value"})

		err := cmd.Execute()

		require.NoError(t, err)
	})
}
