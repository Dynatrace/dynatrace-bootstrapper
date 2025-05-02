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
	var confContent strings.Builder

	if fc.containerSection != nil {
		content, err := fc.containerSection.toString()
		if err != nil {
			return "", err
		}

		confContent.WriteString(content)
		confContent.WriteString("\n")
	}

	if fc.hostSection != nil {
		content, err := fc.hostSection.toString()
		if err != nil {
			return "", err
		}

		confContent.WriteString(content)
		confContent.WriteString("\n")
	}

	return confContent.String(), nil
}

type containerSection struct {
	NodeName                string `json:"k8s_node_name"`
	PodName                 string `json:"k8s_fullpodname"`
	PodUID                  string `json:"k8s_poduid"`
	PodNamespace            string `json:"k8s_namespace"`
	ClusterID               string `json:"k8s_cluster_id"`
	ContainerName           string `json:"k8s_containername"`
	DeprecatedContainerName string `json:"containerName"`
	ImageName               string `json:"imageName"`
}

func (cs containerSection) toMap() (map[string]string, error) {
	return structs.ToMap(cs)
}

func (cs containerSection) toString() (string, error) {
	var sectionContent strings.Builder

	contentMap, err := cs.toMap()
	if err != nil {
		return "", err
	}

	sectionContent.WriteString("[container]")
	sectionContent.WriteString("\n")

	for key, value := range contentMap {
		if value == "" {
			continue
		}

		sectionContent.WriteString(key)
		sectionContent.WriteString(" ")
		sectionContent.WriteString(value)
		sectionContent.WriteString("\n")
	}

	return sectionContent.String(), nil
}

type hostSection struct {
	Tenant      string `json:"tenant"`
	IsFullStack string `json:"isCloudNativeFullStack"`
}

func (hs hostSection) toMap() (map[string]string, error) {
	return structs.ToMap(hs)
}

func (hs hostSection) toString() (string, error) {
	var sectionContent strings.Builder

	contentMap, err := hs.toMap()
	if err != nil {
		return "", err
	}

	sectionContent.WriteString("[host]")
	sectionContent.WriteString("\n")

	for key, value := range contentMap {
		if value == "" {
			continue
		}

		sectionContent.WriteString(key)
		sectionContent.WriteString(" ")
		sectionContent.WriteString(value)
		sectionContent.WriteString("\n")
	}

	return sectionContent.String(), nil
}

func fromAttributes(containerAttr container.Attributes, podAttr pod.Attributes, isFullStack bool) fileContent {
	fileContent := fileContent{
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
		fileContent.hostSection = &hostSection{
			Tenant:      podAttr.DTTenantUID,
			IsFullStack: "true",
		}
		fileContent.NodeName = podAttr.NodeName
	}

	return fileContent
}
