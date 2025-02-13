package enrichment

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	"github.com/spf13/afero"
)

type content struct {
	pod.Attributes `json:",inline"`

	ContainerName string `json:"k8s.container.name"`

	// Deprecated
	DTClusterID string `json:"dt.kubernetes.cluster.id"`
	// Deprecated
	DTWorkloadKind string `json:"dt.kubernetes.workload.kind"`
	// Deprecated
	DTWorkloadName string `json:"dt.kubernetes.workload.name"`
}

func Configure(fs afero.Afero, configDirectory string, podAttr pod.Attributes, containerAttr container.Attributes) error {
	contentJson := content{
		Attributes:     podAttr,
		ContainerName:  containerAttr.ContainerName,
		DTClusterID:    podAttr.ClusterUId,
		DTWorkloadKind: podAttr.WorkloadKind,
		DTWorkloadName: podAttr.WorkloadName,
	}

	raw, err := json.Marshal(data)
	if err != nil {
		return err
	}

	content := map[string]string{}

	err = json.Unmarshal(raw, &content)
	if err != nil {
		return err
	}

	for key, value := range runner.env.WorkloadAnnotations {
		content[key] = value
	}

	jsonContent, err := json.Marshal(content)
	if err != nil {
		return err
	}

	err = runner.createConfigFile(fmt.Sprintf(enrichmentJsonPathTemplate, container.Name), string(jsonContent), true)
	if err != nil {
		return err
	}

	var propsContent strings.Builder
	for key, value := range content {
		propsContent.WriteString(key)
		propsContent.WriteString("=")
		propsContent.WriteString(value)
		propsContent.WriteString("\n")
	}

	err = runner.createConfigFile(fmt.Sprintf(enrichmentPropsPathTemplate, container.Name), propsContent.String(), true)
	if err != nil {
		return err
	}

}
