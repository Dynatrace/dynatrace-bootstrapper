package pod

import (
	"maps"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAttributes(t *testing.T) {
	t.Run("valid attributes", func(t *testing.T) {
		attributes := []string{"k8s.pod.name=pod1", "k8s.pod.uid=123", "k8s.namespace.name=default"}
		expected := Attributes{
			UserDefined: map[string]string{},
			PodInfo: PodInfo{
				PodName:       "pod1",
				PodUID:        "123",
				NamespaceName: "default"},
		}

		result, err := ParseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("empty input => should be ignored", func(t *testing.T) {
		attributes := []string{}
		expected := Attributes{UserDefined: map[string]string{}}
		result, err := ParseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("invalid format => should be ignored", func(t *testing.T) {
		attributes := []string{"invalidEntry"}
		expected := Attributes{UserDefined: map[string]string{}}
		result, err := ParseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("mixed valid and invalid attributes :=> only valid input should be considered", func(t *testing.T) {
		attributes := []string{"k8s.pod.name=pod2", "invalidEntry", "k8s.namespace.name=prod"}
		expected := Attributes{
			UserDefined: map[string]string{},
			PodInfo: PodInfo{
				PodName:       "pod2",
				NamespaceName: "prod",
			},
		}
		result, err := ParseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("mixed valid, invalid  and user-defined attributes :=> only valid and user-defined input should be considered", func(t *testing.T) {
		attributes := []string{"k8s.pod.name=pod2", "invalidEntry", "k8s.namespace.name=prod", "beep=boop"}
		expected := Attributes{
			UserDefined: map[string]string{
				"beep": "boop",
			},
			PodInfo: PodInfo{
				PodName:       "pod2",
				NamespaceName: "prod",
			},
		}
		result, err := ParseAttributes(attributes)
		require.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}

func TestFilterOutUserDefined(t *testing.T) {
	parsedInput := Attributes{
		PodInfo: PodInfo{
			PodName:       "pod2",
			NamespaceName: "prod",
		},
	}

	parseableInput, err := parsedInput.ToMap()
	require.NoError(t, err)

	expectedUserDefined := map[string]string{
		"beep": "boop",
		"tip":  "top",
	}

	rawInput := maps.Clone(parseableInput)
	maps.Copy(rawInput, expectedUserDefined)

	err = filterOutUserDefined(rawInput, parsedInput)
	require.NoError(t, err)
	assert.Equal(t, expectedUserDefined, rawInput)
}

func TestToArgs(t *testing.T) {
	attributes := Attributes{
		UserDefined: map[string]string{
			"beep": "boop",
		},
		PodInfo: PodInfo{
			PodName:       "pod2",
			NamespaceName: "prod",
		},
	}

	expectedArgs := []string{
		"--" + Flag + "=beep=boop",
		"--" + Flag + "=k8s.pod.name=pod2",
		"--" + Flag + "=k8s.namespace.name=prod",
	}

	args, err := ToArgs(attributes)
	require.NoError(t, err)
	assert.ElementsMatch(t, expectedArgs, args)
}
