package metadata

import (
	"encoding/json"
	"maps"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
	"github.com/pkg/errors"
)

type content struct {
	pod.Attributes `json:",inline"`

	ContainerName string `json:"k8s.container.name"`

	// Deprecated
	DTClusterID string `json:"dt.kubernetes.cluster.id,omitempty"`
	// Deprecated
	DTWorkloadKind string `json:"dt.kubernetes.workload.kind,omitempty"`
	// Deprecated
	DTWorkloadName string `json:"dt.kubernetes.workload.name,omitempty"`
}

func (c content) toMap() (map[string]string, error) {
	baseMap, err := structs.ToMap(c)
	if err != nil {
		return nil, err
	}

	maps.Copy(baseMap, c.Attributes.UserDefined)

	return baseMap, nil
}

func (c content) toJson() ([]byte, error) {
	rawMap, err := c.ToMap() // needed to make the pod.Attributes.UserDefined visible
	if err != nil {
		return nil, err
	}

	raw, err := json.Marshal(rawMap)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return raw, nil
}

func (c content) toProperties() (string, error) {
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

func fromAttributes(containerAttr container.Attributes, podAttr pod.Attributes) content {
	return content{
		Attributes:     podAttr,
		ContainerName:  containerAttr.ContainerName,
		DTClusterID:    podAttr.ClusterUId,
		DTWorkloadKind: podAttr.WorkloadKind,
		DTWorkloadName: podAttr.WorkloadName,
	}
}
