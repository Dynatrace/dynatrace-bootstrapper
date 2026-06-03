package k8sinit

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/move"
	"github.com/stretchr/testify/require"
)

func TestBootstrapper(t *testing.T) {
	t.Run("should validate required flags - missing flags -> error", func(t *testing.T) {
		cmd := New()

		err := cmd.Execute()

		require.Error(t, err)
	})
	t.Run("should validate required flags - present flags -> no error", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupSource(t, tmpDir, "123")

		cmd := New()
		cmd.SetArgs([]string{"--source", tmpDir, "--target", t.TempDir()})

		err := cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("--suppress-error=true -> no error", func(t *testing.T) {
		cmd := New()
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\"})

		err := cmd.Execute()

		require.Error(t, err)

		cmd = New()
		// Note: we can't skip/mask the validation of the required flags
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\", "--suppress-error"})

		err = cmd.Execute()

		require.NoError(t, err)

		cmd = New()
		// Note: we can't skip/mask the validation of the required flags
		cmd.SetArgs([]string{"--source", "\\\\", "--target", "\\\\", "--suppress-error", "true"})

		err = cmd.Execute()

		require.NoError(t, err)
	})

	t.Run("should allow unknown flags -> no error", func(t *testing.T) {
		tmpDir := t.TempDir()
		setupSource(t, tmpDir, "123")

		cmd := New()
		cmd.SetArgs([]string{"--source", tmpDir, "--target", t.TempDir(), "--unknown", "--flag", "value"})

		err := cmd.Execute()

		require.NoError(t, err)
	})
}

func TestEnableAttributesDTKubernetes(t *testing.T) {
	const containerName = "test-container"

	podAttrArgs := []string{
		"--attribute=k8s.cluster.uid=test-cluster-uid",
		"--attribute=k8s.workload.kind=Deployment",
		"--attribute=k8s.workload.name=test-workload",
	}
	containerAttrArg := `--attribute-container={"k8s.container.name": "` + containerName + `"}`

	t.Run("flag defaults to true", func(t *testing.T) {
		cmd := New()
		flag := cmd.Flags().Lookup(EnableAttributesDTKubernetesFlag)
		require.NotNil(t, flag)
		require.Equal(t, "true", flag.DefValue)
	})

	t.Run("deprecated dt.kubernetes attributes included when flag is true (default)", func(t *testing.T) {
		srcDir := t.TempDir()
		cfgDir := t.TempDir()
		inputDir := t.TempDir()
		setupSource(t, srcDir, "123")

		cmd := New()
		cmd.SetArgs(append([]string{
			"--source", srcDir,
			"--target", t.TempDir(),
			"--config-directory", cfgDir,
			"--input-directory", inputDir,
			containerAttrArg,
		}, podAttrArgs...))

		require.NoError(t, cmd.Execute())

		jsonContent, err := os.ReadFile(filepath.Join(cfgDir, containerName, "enrichment", "dt_metadata.json"))
		require.NoError(t, err)
		require.Contains(t, string(jsonContent), "dt.kubernetes.cluster.id")
		require.Contains(t, string(jsonContent), "dt.kubernetes.workload.kind")
		require.Contains(t, string(jsonContent), "dt.kubernetes.workload.name")
	})

	t.Run("deprecated dt.kubernetes attributes excluded when flag is false", func(t *testing.T) {
		srcDir := t.TempDir()
		cfgDir := t.TempDir()
		inputDir := t.TempDir()
		setupSource(t, srcDir, "123")

		cmd := New()
		cmd.SetArgs(append([]string{
			"--source", srcDir,
			"--target", t.TempDir(),
			"--config-directory", cfgDir,
			"--input-directory", inputDir,
			"--" + EnableAttributesDTKubernetesFlag + "=false",
			containerAttrArg,
		}, podAttrArgs...))

		require.NoError(t, cmd.Execute())

		jsonContent, err := os.ReadFile(filepath.Join(cfgDir, containerName, "enrichment", "dt_metadata.json"))
		require.NoError(t, err)
		require.NotContains(t, string(jsonContent), "dt.kubernetes.cluster.id")
		require.NotContains(t, string(jsonContent), "dt.kubernetes.workload.kind")
		require.NotContains(t, string(jsonContent), "dt.kubernetes.workload.name")
	})
}

func setupSource(t *testing.T, folder, version string) {
	t.Helper()

	versionFilePath := filepath.Join(folder, move.InstallerVersionFilePath)
	require.NoError(t, os.MkdirAll(filepath.Dir(versionFilePath), os.ModePerm))
	require.NoError(t, os.WriteFile(versionFilePath, []byte(version), 0600))

	agentBinFolder := filepath.Join(folder, filepath.Dir(move.CurrentDir), version)
	require.NoError(t, os.MkdirAll(agentBinFolder, 0700))
}
