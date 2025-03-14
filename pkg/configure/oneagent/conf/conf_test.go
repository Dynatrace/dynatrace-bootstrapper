package conf

import (
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	"github.com/go-logr/zapr"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestConfigure(t *testing.T) {
	podAttr := pod.Attributes{
		PodInfo: pod.PodInfo{
			PodName:       "podname",
			PodUID:        "poduid",
			NamespaceName: "namespacename",
		},
		ClusterInfo: pod.ClusterInfo{
			ClusterUID: "clusteruid",
		},
	}
	containerAttr := container.Attributes{
		ContainerName: "containername",
		ImageInfo: container.ImageInfo{
			Registry:    "registry",
			Repository:  "repository",
			Tag:         "tag",
			ImageDigest: "imagedigest",
		},
	}
	configDir := "path/conf"

	t.Run("success", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		err := Configure(testLog, fs, configDir, containerAttr, podAttr)
		require.NoError(t, err)

		expectedContent, err := fromAttributes(containerAttr, podAttr).toMap()
		require.NoError(t, err)

		content, err := fs.ReadFile(filepath.Join(configDir, configPath))
		require.NoError(t, err)

		for key, value := range expectedContent {
			assert.Contains(t, string(content), key+" "+value)
		}
	})
}
