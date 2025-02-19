package container

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAttributes(t *testing.T) {
	t.Run("valid attributes", func(t *testing.T) {
		attributes = []string{
			`{"container_image.registry": "some.reg.io", "container_image.repository": "test-repo", "container_image.tags": "latest", "container_image.digest": "sha256:abcd1234", "k8s.container.name": "test-container-name"}`,
		}

		expected := []Attributes{
			{
				ImageInfo: ImageInfo{
					Registry:    "some.reg.io",
					Repository:  "test-repo",
					Tag:         "latest",
					ImageDigest: "sha256:abcd1234"},
				ContainerName: "test-container-name",
			},
		}

		result, err := ParseAttributes()
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("empty input => should return empty list", func(t *testing.T) {
		attributes = []string{}
		result, err := ParseAttributes()
		require.NoError(t, err)
		assert.Empty(t, result)
	})
	t.Run("invalid JSON format => should return an error", func(t *testing.T) {
		attributes = []string{"invalid_json"}
		result, err := ParseAttributes()
		require.Error(t, err)
		assert.Nil(t, result)
	})
	t.Run("mixed valid and invalid attributes => only valid input should be considered", func(t *testing.T) {
		attributes = []string{
			`{"container_image.registry":"some.reg.io","container_image.repository":"test-repo","container_image.tags":"latest","container_image.digest":"sha256:abcd1234","k8s.container.name":"test-container"}`,
			"invalid_json",
		}

		expected := []Attributes{
			{
				ImageInfo: ImageInfo{
					Registry:    "some.reg.io",
					Repository:  "test-repo",
					Tag:         "latest",
					ImageDigest: "sha256:abcd1234",
				},
				ContainerName: "test-container-name",
			},
		}

		result, err := ParseAttributes()
		require.Error(t, err)
		assert.NotEqual(t, expected, result)
		assert.Nil(t, result)
	})
}
