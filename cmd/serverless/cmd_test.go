package serverless

import (
	"bytes"
	"log"
	"testing"
	"time"

	"github.com/go-logr/stdr"
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

	t.Run("set 'keep-alive=true' and verify the bootstrapper runs for 5 seconds", func(t *testing.T) {
		started := make(chan struct{})
		finished := make(chan error, 1)

		var buf bytes.Buffer
		stdLogger := log.New(&buf, "", 0)

		// set the standard logger for serverless to enable assertions during testing
		stdr.SetVerbosity(1)
		SetLogger(stdr.New(stdLogger))

		go func() {
			close(started)
			cmd := New()

			cmd.SetArgs([]string{"--keep-alive=true", "--debug", "--target", t.TempDir()})
			finished <- cmd.Execute()
		}()

		// wait for keep-alive goroutine to start
		select {
		case <-started:
			t.Log("keep-alive goroutine has started")
		case <-time.After(5 * time.Second):
			t.Fatal("keep-alive goroutine did not start within 5 seconds")
		}

		// test whether keep-alive is running for 5 seconds
		select {
		case err := <-finished:
			t.Fatalf("the Bootstrapper finished execution in 'keep-alive=true' mode: %v", err)
		case <-time.After(5 * time.Second):
			require.Contains(t, buf.String(), "Running in keep-alive mode")
			t.Log("keep-alive ran for 5 seconds")
		}
	})
}
