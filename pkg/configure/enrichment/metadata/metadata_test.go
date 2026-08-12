package metadata

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/k8sinit/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/k8sinit/configure/attributes/pod"
	"github.com/go-logr/zapr"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

const (
	alwaysEnableDeprecatedAttributes = true
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

		err := Configure(testLog, configDir, podAttr, containerAttr, alwaysEnableDeprecatedAttributes)
		require.NoError(t, err)

		expectedContent, err := fromAttributes(containerAttr, podAttr, alwaysEnableDeprecatedAttributes).toMap()
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

	t.Run("newline injected in user-defined annotation => no property injected", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")

		cursedPodAttr := podAttr
		cursedPodAttr.UserDefined = map[string]string{
			"beep": "boop\ncursed.key=cursed-value",
		}

		err := Configure(testLog, configDir, cursedPodAttr, containerAttr, alwaysEnableDeprecatedAttributes)
		require.NoError(t, err)

		propsContent, err := os.ReadFile(filepath.Join(configDir, PropertiesFilePath))
		require.NoError(t, err)

		for line := range strings.SplitSeq(string(propsContent), "\n") {
			assert.NotEqual(t, "cursed.key=cursed-value", line, "properties must not be sneaked in via control characters")
		}

		assert.Contains(t, string(propsContent), "beep=boopcursed.key=cursed-value")
	})
}
