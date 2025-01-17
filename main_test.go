package main

import (
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestBootstrapper(t *testing.T) {
	t.Run("should validate required flags - missing flags -> error", func(t *testing.T) {
		cmd := bootstrapper(afero.NewMemMapFs())

		err := cmd.Execute()

		require.Error(t, err)
	})
	t.Run("should validate required flags - present flags -> no error", func(t *testing.T) {
		cmd := bootstrapper(afero.NewMemMapFs())
		cmd.SetArgs([]string{"--source", "./", "--target", "./"})

		err := cmd.Execute()

		require.NoError(t, err)
	})
}
