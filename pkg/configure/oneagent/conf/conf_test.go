package conf

import (
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

var testLog = zapr.NewLogger(zap.NewExample())

func TestConfigure(t *testing.T) {
	podAttr := pod.Attributes{
		PodInfo: pod.PodInfo{
			PodName:       "podname",
			PodUID:        "poduid",
			NamespaceName: "namespacename",
			NodeName:      "nodename",
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

	t.Run("success - not fullstack", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")

		err := Configure(testLog, configDir, containerAttr, podAttr, "", false)
		require.NoError(t, err)

		expectedMap, err := fromAttributes(containerAttr, podAttr, "", false).toMap()
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(configDir, ConfigPath))
		require.NoError(t, err)

		missingEntries := []string{}

		for key, value := range expectedMap {
			if value == "" {
				assert.NotContains(t, string(content), key)
				missingEntries = append(missingEntries, key)
			} else {
				assert.Contains(t, string(content), key+" "+value)
			}
		}

		expectedMissingEntries := []string{"tenant", "k8s_node_name", "isCloudNativeFullStack"}
		require.Subset(t, expectedMissingEntries, missingEntries) // incase of isFullstack, the host section is missing from the map

		for _, key := range expectedMissingEntries {
			assert.NotContains(t, string(content), key)
		}

		assert.Contains(t, string(content), "[container]")
		assert.NotContains(t, string(content), "[host]")
	})

	t.Run("success - fullstack", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")

		tenant := "test-tenant"

		err := Configure(testLog, configDir, containerAttr, podAttr, tenant, true)
		require.NoError(t, err)

		expectedMap, err := fromAttributes(containerAttr, podAttr, tenant, true).toMap()
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(configDir, ConfigPath))
		require.NoError(t, err)

		for key, value := range expectedMap {
			assert.Contains(t, string(content), key+" "+value)
		}

		assert.Contains(t, string(content), "[container]")
		assert.Contains(t, string(content), "[host]")
	})

	t.Run("error - fullstack but no tenant", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")

		err := Configure(testLog, configDir, containerAttr, podAttr, "", true)
		require.Error(t, err)
	})

	t.Run("newline injected in attribute => no new section created", func(t *testing.T) {
		baseTempDir := filepath.Join(t.TempDir(), "path")
		configDir := filepath.Join(baseTempDir, "config")

		cursedPodAttr := podAttr
		cursedPodAttr.PodName = "podname\n[host]\ntenant cursed-tenant\nisCloudNativeFullStack lol"

		err := Configure(testLog, configDir, containerAttr, cursedPodAttr, "", false)
		require.NoError(t, err)

		content, err := os.ReadFile(filepath.Join(configDir, ConfigPath))
		require.NoError(t, err)

		for line := range strings.SplitSeq(string(content), "\n") {
			assert.NotEqual(t, "[host]", line, "injected value must not create a new INI section")
			assert.NotEqual(t, "tenant cursed-tenant", line, "injected value must not create a new INI entry")
		}

		assert.Equal(t, 1, strings.Count(string(content), "[container]"))
		assert.Contains(t, string(content), "k8s_fullpodname podname[host]tenant cursed-tenantisCloudNativeFullStack lol")
	})
}
