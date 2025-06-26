// Package pod provides utilities for handling pod attributes.
package pod

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
)

const (
	// Flag is the command-line flag for pod attributes.
	Flag = "attribute"

	// expectedKeyValueParts defines the expected number of parts when splitting a key=value string.
	expectedKeyValueParts = 2
)

// Attributes represents pod attributes including user-defined, pod, workload, and cluster info.
type Attributes struct {
	UserDefined  map[string]string `json:"-"`
	Info         `json:",inline"`
	WorkloadInfo `json:",inline"`
	ClusterInfo  `json:",inline"`
}

// ToMap converts Attributes to a map[string]string.
func (attr Attributes) ToMap() (map[string]string, error) {
	return structs.ToMap(attr)
}

// Info contains pod-related metadata.
type Info struct {
	PodName       string `json:"k8s.pod.name,omitempty"`
	PodUID        string `json:"k8s.pod.uid,omitempty"`
	NodeName      string `json:"k8s.node.name,omitempty"`
	NamespaceName string `json:"k8s.namespace.name,omitempty"`
}

// WorkloadInfo contains workload-related metadata.
type WorkloadInfo struct {
	WorkloadKind string `json:"k8s.workload.kind,omitempty"`
	WorkloadName string `json:"k8s.workload.name,omitempty"`
}

// ClusterInfo contains cluster-related metadata.
type ClusterInfo struct {
	ClusterUID      string `json:"k8s.cluster.uid,omitempty"`
	ClusterName     string `json:"k8s.cluster.name,omitempty"`
	DTClusterEntity string `json:"dt.entity.kubernetes_cluster,omitempty"`
}

// ParseAttributes parses a slice of raw attribute strings into an Attributes struct.
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

	result.UserDefined = rawMap

	return result, nil
}

// ToArgs is a helper func to convert a pod.Attributes to a list of args that can be put into a Pod Template.
func ToArgs(attributes Attributes) ([]string, error) {
	var args = make([]string, 0, len(attributes.UserDefined))

	for k, v := range attributes.UserDefined {
		args = append(args, fmt.Sprintf("--%s=%s", k, v))
	}

	return args, nil
}

// filterOutUserDefined removes known structured fields from the raw map, leaving only user-defined fields.
func filterOutUserDefined(raw map[string]string, _ Attributes) {
	// Remove known structured fields from raw map
	delete(raw, "k8s.pod.name")
	delete(raw, "k8s.pod.uid")
	delete(raw, "k8s.node.name")
	delete(raw, "k8s.namespace.name")
	delete(raw, "k8s.workload.kind")
	delete(raw, "k8s.workload.name")
	delete(raw, "k8s.cluster.uid")
	delete(raw, "k8s.cluster.name")
	delete(raw, "dt.entity.kubernetes_cluster")
}
