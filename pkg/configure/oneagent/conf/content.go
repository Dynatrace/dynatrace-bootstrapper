package conf

import (
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/utils/structs"
)

type content struct {
	PodName                 string `json:"k8s_fullpodname"`
	PodUID                  string `json:"k8s_poduid"`
	PodNamespace            string `json:"k8s_namespace"`
	ClusterID               string `json:"k8s_cluster_id"`
	ContainerName           string `json:"k8s_containername"`
	DeprecatedContainerName string `json:"containerName"`
	ImageName               string `json:"imageName"`
}

func (c content) toMap() (map[string]string, error) {
	return structs.ToMap(c)
}

func (c content) toString() (string, error) {
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

func fromAttributes(containerAttr container.Attributes, podAttr pod.Attributes) content {
	return content{
		PodName:                 podAttr.PodName,
		PodUID:                  podAttr.PodUid,
		PodNamespace:            podAttr.NamespaceName,
		ClusterID:               podAttr.ClusterUId,
		ContainerName:           containerAttr.ContainerName,
		DeprecatedContainerName: containerAttr.ContainerName,
		ImageName:               containerAttr.ImageInfo.ToURI(),
	}
}
