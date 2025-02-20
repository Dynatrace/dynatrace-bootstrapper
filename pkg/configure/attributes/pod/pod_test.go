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
			UserDefined: map[string]string{},
			PodInfo: PodInfo{
				PodName:       "pod1",
				PodUid:        "123",
				NamespaceName: "default"},
		}

		result, err := parseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("empty input => should be ignored", func(t *testing.T) {
		attributes = []string{}
		expected := Attributes{UserDefined: map[string]string{},}
		result, err := parseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("invalid format => should be ignored", func(t *testing.T) {
		attributes = []string{"invalidEntry"}
		expected := Attributes{UserDefined: map[string]string{}}
		result, err := parseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("mixed valid and invalid attributes => only valid input should be considered", func(t *testing.T) {
		attributes = []string{"k8s.pod.name=pod2", "invalidEntry", "k8s.namespace.name=prod"}
		expected := Attributes{
			UserDefined: map[string]string{},
			PodInfo: PodInfo{
				PodName:       "pod2",
				NamespaceName: "prod",
			},
		}
		result, err := parseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("mixed valid, invalid  and user-defined attributes => only valid and user-defined input should be considered", func(t *testing.T) {
		attributes = []string{"k8s.pod.name=pod2", "invalidEntry", "k8s.namespace.name=prod", "beep=boop"}
		expected := Attributes{
			UserDefined: map[string]string{
				"beep": "boop",
			},
			PodInfo: PodInfo{
				PodName:       "pod2",
				NamespaceName: "prod",
			},
		}
		result, err := parseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
