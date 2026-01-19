package serverless

import (
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/deployment"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/tests"
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
		const agentVersion = "1.327.30.20251107-111521"

		sourceDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceDir, agentVersion)

		cmd := New()

		targetDir := t.TempDir()
		cmd.SetArgs([]string{"--source", sourceDir, "--keep-alive=false", "--target", targetDir, "--work", t.TempDir()})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("no error if an unknown parameters are provided", func(t *testing.T) {
		const agentVersion = "1.327.30.20251107-111521"

		sourceDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceDir, agentVersion)

		cmd := New()

		targetDir := t.TempDir()
		cmd.SetArgs([]string{"--source", sourceDir, "--keep-alive=false", "--target", targetDir, "--work", t.TempDir(), "--unknown", "--flag", "value"})

		err := cmd.Execute()
		require.NoError(t, err)
	})

	t.Run("set 'keep-alive=true' and verify the bootstrapper runs for 5 seconds", func(t *testing.T) {
		logsObserver := setupServerlessLogger()
		started := make(chan struct{})
		finished := make(chan error, 1)

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
			tests.RequireLogMessage(t, logsObserver, "Running in keep-alive mode...")
			t.Log("keep-alive ran for 5 seconds")
		}
	})

	t.Run("the deployment status is 'Deployed'", func(t *testing.T) {
		logsObserver := setupServerlessLogger()
		const agentVersion = "1.327.30.20251107-111521"

		sourceDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceDir, agentVersion)

		targetDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetDir, agentVersion, agentVersion)

		cmd := New()
		cmd.SetArgs([]string{"--keep-alive=false", "--source", sourceDir, "--target", targetDir})
		err := cmd.Execute()

		tests.RequireLogMessage(t, logsObserver, "OneAgent is already deployed")
		require.NoError(t, err)
	})

	t.Run("the deployment status is 'Not Deployed'", func(t *testing.T) {
		logsObserver := setupServerlessLogger()
		const agentVersion = "1.327.30.20251107-111521"

		sourceDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceDir, agentVersion)

		targetDir := t.TempDir()

		cmd := New()
		cmd.SetArgs([]string{"--keep-alive=false", "--source", sourceDir, "--target", targetDir, "--work", t.TempDir()})
		err := cmd.Execute()
		require.NoError(t, err)

		tests.RequireLogMessage(t, logsObserver, "OneAgent deployment status", "status", "Not deployed")
	})

	t.Run("the deployment status is 'Link Missing'", func(t *testing.T) {
		logsObserver := setupServerlessLogger()
		const agentVersion = "1.327.30.20251107-111521"

		sourceDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceDir, agentVersion)

		targetDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetDir, agentVersion, "")

		cmd := New()
		cmd.SetArgs([]string{"--keep-alive=false", "--source", sourceDir, "--target", targetDir, "--work", t.TempDir()})
		err := cmd.Execute()
		require.NoError(t, err)

		tests.RequireLogMessage(t, logsObserver, "OneAgent deployment status", "status", "Deployment is not complete")
	})

	t.Run("the deployment status returns an error", func(t *testing.T) {
		logsObserver := setupServerlessLogger()

		const agentVersion = "1.327.30.20251107-111521"

		sourceDir := t.TempDir()
		tests.SetupSourceDirectory(t, sourceDir, agentVersion)

		targetDir := t.TempDir()
		tests.SetupTargetDirectory(t, targetDir, agentVersion, agentVersion)

		// change the permission of the target folder, so the deployment check must result into an error
		agentTargetDirectory := deployment.GetAgentFolder(targetDir, agentVersion)
		parentDir := filepath.Dir(agentTargetDirectory)
		err := os.Chmod(parentDir, 0000)
		defer func() {
			require.NoError(t, os.Chmod(parentDir, 0700)) // restore permissions on exit to allow cleanup of the temporary directory
		}()

		require.NoError(t, err)

		cmd := New()
		cmd.SetArgs([]string{"--keep-alive=false", "--source", sourceDir, "--target", targetDir})
		err = cmd.Execute()
		require.ErrorIs(t, err, syscall.EACCES)
		require.Regexp(t, "cannot obtain OneAgent directory info: stat .+: permission denied", err)

		tests.RequireLogMessage(t, logsObserver, "failed to check OneAgent deployment status. Skipping deployment.", "status", "Unknown")
	})
}

// setupServerlessLogger sets the test logger as the default Serverless logger
// and returns a CapturedLogs instance to be used in tests for log message assertions.
func setupServerlessLogger() *tests.CapturedLogs {
	log, capturedLogs := tests.NewTestLogger()
	SetLogger(log)

	return capturedLogs
}
