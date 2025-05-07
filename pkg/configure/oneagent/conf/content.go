package conf

import (
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/pod"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
)

type fileContent struct {
	*containerSection `json:",inline,omitempty"`
	*hostSection      `json:",inline,omitempty"`
}

func (fc fileContent) toMap() (map[string]string, error) {
	return structs.ToMap(fc)
}

func (fc fileContent) toString() (string, error) {
	var content strings.Builder

	if fc.containerSection != nil {
		sectionContent, err := fc.containerSection.toString()
		if err != nil {
			return "", err
		}

		content.WriteString(sectionContent)
		content.WriteString("\n")
	}

	if fc.hostSection != nil {
		sectionContent, err := fc.hostSection.toString()
		if err != nil {
			return "", err
		}

		content.WriteString(sectionContent)
		content.WriteString("\n")
	}

	return content.String(), nil
}

type containerSection struct {
	NodeName                string `json:"k8s_node_name,omitempty"`
	PodName                 string `json:"k8s_fullpodname,omitempty"`
	PodUID                  string `json:"k8s_poduid,omitempty"`
	PodNamespace            string `json:"k8s_namespace,omitempty"`
	ClusterID               string `json:"k8s_cluster_id,omitempty"`
	ContainerName           string `json:"k8s_containername,omitempty"`
	DeprecatedContainerName string `json:"containerName,omitempty"`
	ImageName               string `json:"imageName,omitempty"`
}

func (cs containerSection) toMap() (map[string]string, error) {
	return structs.ToMap(cs)
}

func (cs containerSection) toString() (string, error) {
	var content strings.Builder

	contentMap, err := cs.toMap()
	if err != nil {
		return "", err
	}

	content.WriteString("[container]")
	content.WriteString("\n")

	for key, value := range contentMap {
		if value == "" {
			continue
		}

		content.WriteString(key)
		content.WriteString(" ")
		content.WriteString(value)
		content.WriteString("\n")
	}

	return content.String(), nil
}

type hostSection struct {
	Tenant      string `json:"tenant,omitempty"`
	IsFullStack string `json:"isCloudNativeFullStack,omitempty"`
}

func (hs hostSection) toMap() (map[string]string, error) {
	return structs.ToMap(hs)
}

func (hs hostSection) toString() (string, error) {
	var content strings.Builder

	contentMap, err := hs.toMap()
	if err != nil {
		return "", err
	}

	content.WriteString("[host]")
	content.WriteString("\n")

	for key, value := range contentMap {
		if value == "" {
			continue
		}

		content.WriteString(key)
		content.WriteString(" ")
		content.WriteString(value)
		content.WriteString("\n")
	}

	return content.String(), nil
}

func fromAttributes(containerAttr container.Attributes, podAttr pod.Attributes, tenant string, isFullStack bool) fileContent {
	fc := fileContent{
		containerSection: &containerSection{
			PodName:                 podAttr.PodName,
			PodUID:                  podAttr.PodUID,
			PodNamespace:            podAttr.NamespaceName,
			ClusterID:               podAttr.ClusterUID,
			ContainerName:           containerAttr.ContainerName,
			DeprecatedContainerName: containerAttr.ContainerName,
			ImageName:               containerAttr.ToURI(),
		},
	}

	if isFullStack {
		fc.hostSection = &hostSection{
			Tenant:      tenant,
			IsFullStack: "true",
		}
		fc.NodeName = podAttr.NodeName
	}

	return fc
}
