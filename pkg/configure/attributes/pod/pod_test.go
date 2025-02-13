package pod

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAttributes(t *testing.T) {
	input := []string{
		"k8s.pod.name=random",
		"k8s.pod.uid=random",
		"k8s.namespace.name=random",
		"k8s.cluster.uid=random",
		"k8s.workload.kind=random",
		"k8s.workload.name=random",
		"k8s.container.name=random",
		"beep=boop",
	}
	attr, err := parseAttributes(input)
	require.NoError(t, err)
	assert.NotEmpty(t, attr)

}
