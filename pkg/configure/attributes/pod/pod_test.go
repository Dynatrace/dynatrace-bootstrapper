package pod

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAttributes(t *testing.T) {
	t.Run("valid attributes", func(t *testing.T) {
		attributes = []string{"k8s.pod.name=pod1", "k8s.pod.uid=123", "k8s.namespace.name=default"}
		expected := Attributes{
			PodInfo: PodInfo{
				PodName:       "pod1",
				PodUid:        "123",
				NamespaceName: "default"},
			Raw: map[string]string{
				"k8s.pod.name":       "pod1",
				"k8s.pod.uid":        "123",
				"k8s.namespace.name": "default",
			},
		}

		result, err := ParseAttributes()
		require.NoError(t, err)
		assert.Equal(t, expected, *result)
	})
	t.Run("empty input => should be ignored", func(t *testing.T) {
		attributes = []string{}
		expected := Attributes{Raw: map[string]string{}}
		result, err := ParseAttributes()
		require.NoError(t, err)
		assert.Equal(t, expected, *result)
	})
	t.Run("invalid format => should be ignored", func(t *testing.T) {
		attributes = []string{"invalidEntry"}
		expected := Attributes{Raw: map[string]string{}}
		result, err := ParseAttributes()
		require.NoError(t, err)
		assert.Equal(t, expected, *result)
	})
	t.Run("mixed valid and invalid attributes => only valid input should be considered", func(t *testing.T) {
		attributes = []string{"k8s.pod.name=pod2", "invalidEntry", "k8s.namespace.name=prod"}
		expected := Attributes{
			PodInfo: PodInfo{
				PodName:       "pod2",
				NamespaceName: "prod",
			},
			Raw: map[string]string{
				"k8s.pod.name":       "pod2",
				"k8s.namespace.name": "prod",
			},
		}
		result, err := ParseAttributes()
		require.NoError(t, err)
		assert.Equal(t, expected, *result)
	})
}
