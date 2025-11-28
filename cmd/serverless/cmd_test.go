package serverless

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServerlessCmd(t *testing.T) {
	t.Run("missing required parameters results in an error", func(t *testing.T) {
		cmd := New()
		err := cmd.Execute()

		require.Error(t, err)
		require.ErrorContains(t, err, "required flag(s) \"keep-alive\", \"target\" not set")
	})

	t.Run("missing 'target' parameter results in an error", func(t *testing.T) {
		cmd := New()
		cmd.SetArgs([]string{"--keep-alive=true"})

		err := cmd.Execute()

		require.Error(t, err)
		require.ErrorContains(t, err, "required flag(s) \"target\" not set")
	})

	t.Run("missing 'keep-alive' parameter results in an error", func(t *testing.T) {
		cmd := New()

		targetDir := t.TempDir()
		cmd.SetArgs([]string{"--target", targetDir})

		err := cmd.Execute()

		require.Error(t, err)
		require.ErrorContains(t, err, "required flag(s) \"keep-alive\" not set")
	})

	t.Run("no error if all required parameters are provided", func(t *testing.T) {
		cmd := New()

		targetDir := t.TempDir()
		cmd.SetArgs([]string{"--keep-alive=false", "--target", targetDir})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("no error if an unknown parameters are provided", func(t *testing.T) {
		cmd := New()

		targetDir := t.TempDir()
		cmd.SetArgs([]string{"--keep-alive=false", "--target", targetDir, "--unknown", "--flag", "value"})

		err := cmd.Execute()
		require.NoError(t, err)
	})
}
