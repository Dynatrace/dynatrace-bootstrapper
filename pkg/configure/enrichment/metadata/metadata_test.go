package metadata

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/pod"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

var testLog = zapr.NewLogger(zap.NewExample())

func TestConfigure(t *testing.T) {
	podAttr := pod.Attributes{
		UserDefined: map[string]string{
			"beep": "boop",
			"tip":  "top",
		},
		PodInfo: pod.PodInfo{
			PodName:       "podname",
			PodUID:        "poduid",
			NodeName:      "nodename",
			NamespaceName: "namespacename",
		},
		ClusterInfo: pod.ClusterInfo{
			ClusterUID:      "clusteruid",
			ClusterName:     "clustername",
			DTClusterEntity: "dtclusterentity",
		},
		WorkloadInfo: pod.WorkloadInfo{
			WorkloadKind: "workloadkind",
			WorkloadName: "workloadname",
		},
	}
	containerAttr := container.Attributes{
		ContainerName: "containername",
	}

	t.Run("success", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")

		err := Configure(testLog, configDir, podAttr, containerAttr)
		require.NoError(t, err)

		expectedContent, err := fromAttributes(containerAttr, podAttr).toMap()
		require.NoError(t, err)

		jsonFilePath := filepath.Join(configDir, JSONFilePath)
		jsonContent, err := os.ReadFile(jsonFilePath)
		require.NoError(t, err)

		for key, value := range expectedContent {
			assert.Contains(t, string(jsonContent), fmt.Sprintf("\"%s\":\"%s\"", key, value))
		}

		propsContent, err := os.ReadFile(filepath.Join(configDir, PropertiesFilePath))
		require.NoError(t, err)

		for key, value := range expectedContent {
			assert.Contains(t, string(propsContent), key+"="+value)
		}
	})
}
