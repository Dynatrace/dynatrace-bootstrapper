package pod

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
)

const (
	Flag = "attribute"

	expectedKeyValueParts = 2
)

type Attributes struct {
	UserDefined  map[string]string `json:"-"`
	PodInfo      `json:",inline"`
	WorkloadInfo `json:",inline"`
	ClusterInfo  `json:",inline"`
}

func (attr Attributes) ToMap() (map[string]string, error) {
	return structs.ToMap(attr)
}

type PodInfo struct { //nolint:revive
	PodName       string `json:"k8s.pod.name,omitempty"`
	PodUID        string `json:"k8s.pod.uid,omitempty"`
	NodeName      string `json:"k8s.node.name,omitempty"`
	NamespaceName string `json:"k8s.namespace.name,omitempty"`
}

type WorkloadInfo struct {
	WorkloadKind string `json:"k8s.workload.kind,omitempty"`
	WorkloadName string `json:"k8s.workload.name,omitempty"`
}

type ClusterInfo struct {
	ClusterUID      string `json:"k8s.cluster.uid,omitempty"`
	ClusterName     string `json:"k8s.cluster.name,omitempty"`
	DTClusterEntity string `json:"dt.entity.kubernetes_cluster,omitempty"`
}

func ParseAttributes(rawAttributes []string) (Attributes, error) {
	rawMap := make(map[string]string, len(rawAttributes))

	for _, attr := range rawAttributes {
		parts := strings.Split(attr, "=")
		if len(parts) == expectedKeyValueParts {
			rawMap[parts[0]] = parts[1]
		}
	}

	raw, err := json.Marshal(rawMap)
	if err != nil {
		return Attributes{}, err
	}

	var result Attributes

	err = json.Unmarshal(raw, &result)
	if err != nil {
		return Attributes{}, err
	}

	err = filterOutUserDefined(rawMap, result)
	if err != nil {
		return Attributes{}, err
	}

	result.UserDefined = rawMap

	return result, nil
}

func filterOutUserDefined(rawInput map[string]string, parsedInput Attributes) error {
	parsedMap, err := parsedInput.ToMap()
	if err != nil {
		return err
	}

	for key := range parsedMap {
		delete(rawInput, key)
	}

	return nil
}

// ToArgs is a helper func to convert an pod.Attributes to a list of args that can be put into a Pod Template
func ToArgs(attributes Attributes) ([]string, error) {
	attrMap, err := attributes.ToMap()
	if err != nil {
		return nil, err
	}

	args := make([]string, 0, len(attrMap))

	for key, value := range attrMap {
		if value == "" {
			continue
		}

		args = append(args, fmt.Sprintf("--%s=%s=%s", Flag, key, value))
	}

	for key, value := range attributes.UserDefined {
		if value == "" {
			continue
		}

		args = append(args, fmt.Sprintf("--%s=%s=%s", Flag, key, value))
	}

	return args, nil
}
