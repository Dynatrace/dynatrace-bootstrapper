package config_files

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/container"
	"github.com/Dynatrace/dynatrace-bootstrapper/pkg/configure/attributes/pod"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	onlyReadAllFileMode      = 0444
	enrichmentJsonFile       = "dt_metadata.json"
	enrichmentPropertiesFile = "dt_metadata.properties"
	containerConfFile        = "container.conf"
)

var (
	oneAgentDir       = filepath.Join("oneagent", "agent", "config")
	enrichmentDir     = "enrichment"
	containerNameKeys = []string{"containerName", "k8s_containername"}
	imageNameKey      = "imageName"
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

func ConfigureEnrichmentFiles(fs afero.Afero, configDirectory string, podAttr pod.Attributes, containerName string) error {
	contentJson := content{
		Attributes:    podAttr,
		ContainerName: containerName,
	}

	if podAttr.ClusterUId != "" {
		contentJson.Raw["dt.kubernetes.cluster.id"] = podAttr.ClusterUId
	}

	if podAttr.WorkloadKind != "" {
		contentJson.Raw["dt.kubernetes.workload.kind"] = podAttr.WorkloadKind
	}

	if podAttr.WorkloadName != "" {
		contentJson.Raw["dt.kubernetes.workload.name"] = podAttr.WorkloadName
	}

	if containerName != "" {
		contentJson.Raw["k8s.container.name"] = containerName
	}

	logrus.Infof("Format content into a raw form: %s", contentJson)

	raw, err := json.Marshal(contentJson.Raw)
	if err != nil {
		logrus.Errorf("Error marshalling content: %s", contentJson)
		return err
	}

	logrus.Infof("Created raw content: %s", raw)

	content := map[string]string{}

	err = json.Unmarshal(raw, &content)
	if err != nil {
		logrus.Errorf("Error unmarshalling content: %s", content)

		return err
	}

	err = createConfigFile(fs, filepath.Join(configDirectory, contentJson.ContainerName, enrichmentDir, enrichmentJsonFile), string(raw))
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

	err = createConfigFile(fs, filepath.Join(configDirectory, contentJson.ContainerName, enrichmentDir, enrichmentPropertiesFile), propsContent.String())
	if err != nil {
		return err
	}

	return nil
}

func ConfigureContainerConfFile(fs afero.Afero, configDirectory string, containerAttr container.Attributes) error {
	raw, err := json.Marshal(containerAttr)
	if err != nil {
		logrus.Errorf("Error marshalling content: %s", containerAttr)
		return err
	}

	logrus.Infof("Created raw content: %s", raw)

	content := map[string]string{}

	err = json.Unmarshal(raw, &content)
	if err != nil {
		logrus.Errorf("Error unmarshalling content: %s", content)

		return err
	}

	var containerConfContent strings.Builder

	prepareContainerConfContent(&containerConfContent, content, containerAttr.ImageInfo)

	err = createConfigFile(fs, filepath.Join(configDirectory, containerAttr.ContainerName, oneAgentDir, containerConfFile), containerConfContent.String())
	if err != nil {
		return err
	}

	return nil
}

func prepareContainerConfContent(containerConfContent *strings.Builder, content map[string]string, imageInfo container.ImageInfo) *strings.Builder {
	for key, value := range content {
		if key == "k8s.container.name" {
			containerConfContent.WriteString(key)
			containerConfContent.WriteString(" ")
			containerConfContent.WriteString(value)
			containerConfContent.WriteString("\n")
			addAdditionalContainerNameKeys(containerConfContent, value)
		}
	}

	containerConfContent.WriteString(imageNameKey)
	containerConfContent.WriteString(" ")
	containerConfContent.WriteString(generateImageName(imageInfo))
	containerConfContent.WriteString("\n")

	return containerConfContent
}

func addAdditionalContainerNameKeys(containerConfContent *strings.Builder, value string) *strings.Builder {
	for _, containerNameKey := range containerNameKeys {
		containerConfContent.WriteString(containerNameKey)
		containerConfContent.WriteString(" ")
		containerConfContent.WriteString(value)
		containerConfContent.WriteString("\n")
	}

	return containerConfContent
}

func createConfigFile(fs afero.Afero, path string, content string) error {
	err := createFile(fs, path, content)
	if err != nil {
		return errors.WithStack(err)
	}

	logrus.Infof("Created file at: %s with following content: %s", path, content)

	return nil
}

func createFile(fs afero.Fs, path string, content string) error {
	err := fs.MkdirAll(filepath.Dir(path), onlyReadAllFileMode)
	if err != nil {
		return errors.WithStack(err)
	}

	file, err := fs.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, onlyReadAllFileMode)
	if err != nil {
		return errors.WithStack(err)
	}

	_, err = file.Write([]byte(content))
	if err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func generateImageName(imageInfo container.ImageInfo) string {
	var imageName string

	registry := imageInfo.Registry
	repository := imageInfo.Repository
	tag := imageInfo.Tag
	digest := imageInfo.ImageDigest

	if registry != "" {
		if repository != "" {
			imageName = registry + "/" + repository
		} else {
			imageName = repository
		}

		if digest != "" {
			imageName += "@" + digest
		} else if tag != "" {
			imageName += ":" + tag
		}
	}

	return imageName
}
