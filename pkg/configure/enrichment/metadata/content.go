// Package metadata provides types and helpers for enrichment metadata content.
package metadata

import (
	"encoding/json"
	"maps"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/pod"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
	"github.com/pkg/errors"
)

// fileContent represents the structure of the metadata file content.
type fileContent struct {
	pod.Attributes `json:",inline"`

	ContainerName string `json:"k8s.container.name"`

	// Deprecated
	DTClusterID string `json:"dt.kubernetes.cluster.id,omitempty"`
	// Deprecated
	DTWorkloadKind string `json:"dt.kubernetes.workload.kind,omitempty"`
	// Deprecated
	DTWorkloadName string `json:"dt.kubernetes.workload.name,omitempty"`
}

// toMap converts fileContent to a map[string]string, merging user-defined attributes.
func (c fileContent) toMap() (map[string]string, error) {
	baseMap, err := structs.ToMap(c)
	if err != nil {
		return nil, err
	}

	maps.Copy(baseMap, c.UserDefined)

	return baseMap, nil
}

// toJSON marshals fileContent to JSON, making UserDefined visible.
func (c fileContent) toJSON() ([]byte, error) {
	rawMap, err := c.toMap() // needed to make the pod.Attributes.UserDefined visible
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(rawMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return raw, nil
}

// toProperties converts fileContent to a properties-like string format.
func (c fileContent) toProperties() (string, error) {
	var confContent strings.Builder

	contentMap, err := c.toMap()
	if err != nil {
		return "", err
	}

	for key, value := range contentMap {
		confContent.WriteString(key)
		confContent.WriteString("=")
		confContent.WriteString(value)
		confContent.WriteString("\n")
	}

	return confContent.String(), nil
}

func fromAttributes(containerAttr container.Attributes, podAttr pod.Attributes) fileContent {
	return fileContent{
		Attributes:     podAttr,
		ContainerName:  containerAttr.ContainerName,
		DTClusterID:    podAttr.ClusterUID,
		DTWorkloadKind: podAttr.WorkloadKind,
		DTWorkloadName: podAttr.WorkloadName,
	}
}
