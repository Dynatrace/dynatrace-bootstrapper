package move

import (
	"encoding/json"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteMetadataJsonFile(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	data := map[string]string{
		k8sPodNameKey:       "test-pod-name",
		k8sClusterUidKey:    "test-cluster-uid",
		k8sContainerNameKey: "test-container-name",
		"k8s.node.name":     "test-node-name",
		"k8s.workload.kind": "some-workload",
	}

	expectedJson := `{
	"dt.kubernetes.cluster.id":"test-cluster-uid",
	"dt.kubernetes.workload.kind":"some-workload",
	"k8s.pod.name": "test-pod-name",
	"k8s.cluster.uid": "test-cluster-uid",
	"k8s.container.name": "test-container-name",
	"k8s.node.name" : "test-node-name",
	"k8s.workload.kind": "some-workload"
}`
	filePath := "dt_metadata.json"
	err := writeMetadataJsonFile(fs, filePath, data)
	assert.NoError(t, err)

	content, err := afero.ReadFile(fs, filePath)
	assert.NoError(t, err)

	var expected, actual map[string]string
	err = json.Unmarshal([]byte(expectedJson), &expected)
	require.NoError(t, err)

	err = json.Unmarshal(content, &actual)
	require.NoError(t, err)

	assert.Equal(t, expected, actual)
}

func TestWriteMetadataPropertiesFile(t *testing.T) {
	fs := afero.Afero{Fs: afero.NewMemMapFs()}

	data := map[string]string{
		k8sPodNameKey:    "test-pod-name",
		k8sClusterUidKey: "test-cluster-uid",
		"k8s.node.name":  "test-node-name",
	}

	expectedProperties := "dt.kubernetes.cluster.id test-cluster-uid\nk8s.cluster.uid test-cluster-uid\nk8s.node.name test-node-name\nk8s.pod.name test-pod-name\n"

	filePath := "dt_metadata.properties"
	err := writeMetadataPropertiesFile(fs, filePath, data)
	assert.NoError(t, err)

	content, err := afero.ReadFile(fs, filePath)
	assert.NoError(t, err)
	assert.Equal(t, expectedProperties, string(content))
}

func TestWriteContainerConfFile(t *testing.T) {
	t.Run("container conf with image with tag", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		data := map[string]string{
			k8sContainerNameKey:         "test-container-name",
			containerImageRegistryKey:   "test.registry.io",
			containerImageRepositoryKey: "test-repo",
			containerImageTagsKey:       "v1",
		}

		expectedConfig := "containerName test-container-name\nk8s_containername test-container-name\nimageName test.registry.io/test-repo:v1\n"

		filePath := "container.conf"
		err := writeContainerConfFile(fs, filePath, data)
		assert.NoError(t, err)

		content, err := afero.ReadFile(fs, filePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedConfig, string(content))
	})
	t.Run("container conf with image with digest", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		data := map[string]string{
			k8sContainerNameKey:         "test-container-name",
			containerImageRegistryKey:   "test.registry.io",
			containerImageRepositoryKey: "test-repo",
			containerImageDigestKey:     "sha256:1234",
		}

		expectedConfig := "containerName test-container-name\nk8s_containername test-container-name\nimageName test.registry.io/test-repo@sha256:1234\n"

		filePath := "container.conf"
		err := writeContainerConfFile(fs, filePath, data)
		assert.NoError(t, err)

		content, err := afero.ReadFile(fs, filePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedConfig, string(content))
	})
	t.Run("container conf with image with digest and tag -> digest will be set", func(t *testing.T) {
		fs := afero.Afero{Fs: afero.NewMemMapFs()}

		data := map[string]string{
			k8sContainerNameKey:         "test-container-name",
			containerImageRegistryKey:   "test.registry.io",
			containerImageRepositoryKey: "test-repo",
			containerImageTagsKey:       "v1",
			containerImageDigestKey:     "sha256:1234",
		}

		expectedConfig := "containerName test-container-name\nk8s_containername test-container-name\nimageName test.registry.io/test-repo@sha256:1234\n"

		filePath := "container.conf"
		err := writeContainerConfFile(fs, filePath, data)
		assert.NoError(t, err)

		content, err := afero.ReadFile(fs, filePath)
		assert.NoError(t, err)
		assert.Equal(t, expectedConfig, string(content))
	})
}
