package conf

import (
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/cmd/configure/attributes/pod"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
)

type fileContent struct {
	PodName                 string `json:"k8s_fullpodname"`
	PodUID                  string `json:"k8s_poduid"`
	PodNamespace            string `json:"k8s_namespace"`
	ClusterID               string `json:"k8s_cluster_id"`
	ContainerName           string `json:"k8s_containername"`
	DeprecatedContainerName string `json:"containerName"`
	ImageName               string `json:"imageName"`
}

func (c fileContent) toMap() (map[string]string, error) {
	return structs.ToMap(c)
}

func (c fileContent) toString() (string, error) {
	var confContent strings.Builder

	contentMap, err := c.toMap()
	if err != nil {
		return "", err
	}

	confContent.WriteString("[container]")
	confContent.WriteString("\n")

	for key, value := range contentMap {
		confContent.WriteString(key)
		confContent.WriteString(" ")
		confContent.WriteString(value)
		confContent.WriteString("\n")
	}

	return confContent.String(), nil
}

func fromAttributes(containerAttr container.Attributes, podAttr pod.Attributes) fileContent {
	return fileContent{
		PodName:                 podAttr.PodName,
		PodUID:                  podAttr.PodUID,
		PodNamespace:            podAttr.NamespaceName,
		ClusterID:               podAttr.ClusterUID,
		ContainerName:           containerAttr.ContainerName,
		DeprecatedContainerName: containerAttr.ContainerName,
		ImageName:               containerAttr.ImageInfo.ToURI(),
	}
}
