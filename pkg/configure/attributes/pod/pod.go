package pod

import (
	"encoding/json"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
	"github.com/spf13/cobra"
)

const (
	flagKey = "attribute"
)

var (
	attributes []string
)

func AddFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringArrayVar(&attributes, flagKey, []string{}, "(Optional) Pod-specific attributes in key=value format.")
}

type Attributes struct {
	UserDefined  map[string]string `json:"-"`
	PodInfo      `json:",inline"`
	WorkloadInfo `json:",inline"`
	ClusterInfo  `json:",inline"`
}

func (attr Attributes) ToMap() (map[string]string, error) {
	return structs.ToMap(attr)
}

type PodInfo struct {
	PodName       string `json:"k8s.pod.name"`
	PodUid        string `json:"k8s.pod.uid"`
	NamespaceName string `json:"k8s.namespace.name"`
}

type WorkloadInfo struct {
	WorkloadKind string `json:"k8s.workload.kind"`
	WorkloadName string `json:"k8s.workload.name"`
}

type ClusterInfo struct {
	ClusterUId      string `json:"k8s.cluster.uid"`
	DTClusterEntity string `json:"dt.entity.kubernetes_cluster"`
}

func ParseAttributes() (Attributes, error) {
	return parseAttributes(attributes)
}

func parseAttributes(rawAttributes []string) (Attributes, error) {
	rawMap := make(map[string]string)

	for _, attr := range rawAttributes {
		parts := strings.Split(attr, "=")
		if len(parts) == 2 {
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