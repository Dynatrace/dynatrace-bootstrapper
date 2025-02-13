package move

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
)

const (
	enrichmentDir          = "enrichment"
	oneAgentDir            = "oneagent/agent/config"
	containerConfigFile    = "container.conf"
	metadataJsonFile       = "dt_metadata.json"
	metadataPropertiesFile = "dt_metadata.properties"

	// all key that come into the metadata files
	k8sPodNameKey              = "k8s.pod.name"
	k8sPodUidKey               = "k8s.pod.uid"
	k8sNamespaceNameKey        = "k8s.namespace.name"
	k8sClusterUidKey           = "k8s.cluster.uid"
	k8sWorkloadKindKey         = "k8s.workload.kind"
	k8sWorkloadNameKey         = "k8s.workload.name"
	k8sContainerNameKey        = "k8s.container.name"
	deprecatedClusterUidKey    = "dt.kubernetes.cluster.id"
	depcrecatedWorkloadKindKey = "dt.kubernetes.workload.kind"
	depcrecatedWorkloadNameKey = "dt.kubernetes.workload.name"

	// all key that come into the container config file
	containerConfPodNameKey          = "k8s_fullpodname"
	containerConfPodUidKey           = "k8s_poduid"
	containerConfNamespaceNameKey    = "k8s_namespace"
	containerConfClusterUidKey       = "k8s_cluster_id"
	containerConfContainerNameKey    = "containerName"
	containerConfK8sContainerNameKey = "k8s_containername"

	containerImageRegistryKey   = "container_image.registry"
	containerImageRepositoryKey = "container_image.repository"
	containerImageTagsKey       = "container_image.tags"
	containerImageDigestKey     = "container_image.digest"
)

type AttributeMapping struct {
	SourceKey                 string
	MetadataJsonKey           string
	DeprecatedMetadataJsonKey string
	ContainerConfKeys         []string
	IsImageField              bool
}

var attributeMappings = []AttributeMapping{
	{k8sPodNameKey, k8sPodNameKey, "", []string{containerConfPodNameKey}, false},
	{k8sPodUidKey, k8sPodUidKey, "", []string{containerConfPodUidKey}, false},
	{k8sNamespaceNameKey, k8sNamespaceNameKey, "", []string{containerConfNamespaceNameKey}, false},
	{k8sClusterUidKey, k8sClusterUidKey, deprecatedClusterUidKey, []string{containerConfClusterUidKey}, false},
	{k8sWorkloadKindKey, k8sWorkloadKindKey, depcrecatedWorkloadKindKey, nil, false},
	{k8sWorkloadNameKey, k8sWorkloadNameKey, depcrecatedWorkloadNameKey, nil, false},
	{k8sContainerNameKey, k8sContainerNameKey, "", []string{containerConfContainerNameKey, containerConfK8sContainerNameKey}, false},
	{containerImageRegistryKey, "", "", nil, true},
	{containerImageRepositoryKey, "", "", nil, true},
	{containerImageTagsKey, "", "", nil, true},
	{containerImageDigestKey, "", "", nil, true},
}

func writeConfigFiles(fs afero.Afero) error {
	logrus.Infof("Starting to write config files to: %s", configDirectory)

	podMetadata := parseAttributes(attribute)

	for _, containerAttribute := range attributeContainers {
		containerData, err := parseContainerAttributes(containerAttribute)
		if err != nil {
			return err
		}

		if containerData == nil {
			logrus.Infof("No containerData for %s, skipping", containerAttribute)
			continue
		}

		containerName, exists := containerData[k8sContainerNameKey]
		if !exists {
			logrus.Infof("No containerName for %s, skipping", containerData)
			continue
		}

		containerPath := filepath.Join(targetFolder, configDirectory, containerName)
		enrichmentPath := filepath.Join(containerPath, enrichmentDir)
		agentConfigPath := filepath.Join(containerPath, oneAgentDir)

		for _, dir := range []string{enrichmentPath, agentConfigPath} {
			if err := fs.MkdirAll(dir, 0755); err != nil {
				return err
			}
		}

		metadataFile := filepath.Join(enrichmentPath, metadataJsonFile)
		mergedMetadata := mergeMetadata(podMetadata, containerData)

		if err := writeMetadataJsonFile(fs, metadataFile, mergedMetadata); err != nil {
			return err
		}

		propertiesFile := filepath.Join(enrichmentPath, metadataPropertiesFile)
		if err := writeMetadataPropertiesFile(fs, propertiesFile, mergedMetadata); err != nil {
			return err
		}

		confFile := filepath.Join(agentConfigPath, containerConfigFile)
		if err := writeContainerConfFile(fs, confFile, containerData); err != nil {
			return err
		}
	}

	return nil
}

func parseAttributes(attrs []string) map[string]string {
	logrus.Infof("Starting to parse pod attributes for: %s", attrs)

	result := make(map[string]string)

	for _, attr := range attrs {
		parts := strings.Split(attr, "=")
		if len(parts) == 2 {
			result[parts[0]] = parts[1]
		}
	}

	return result
}

func parseContainerAttributes(jsonStr string) (map[string]string, error) {
	logrus.Infof("Starting to parse container attributes for: %s", jsonStr)

	var attributes map[string]string
	if err := json.Unmarshal([]byte(jsonStr), &attributes); err != nil {
		logrus.Warnf("json format is invalid for: %s", jsonStr)
		return nil, err
	}

	return attributes, nil
}

func mergeMetadata(podAttributes, containerAttributes map[string]string) map[string]string {
	mergedAttributes := make(map[string]string)

	for k, v := range podAttributes {
		mergedAttributes[k] = v
	}

	// since k8s.container.name is also written to dt_metadata files we need to merge it here to the pod attributes
	if containerName, exists := containerAttributes[k8sContainerNameKey]; exists {
		mergedAttributes[k8sContainerNameKey] = containerName
	}

	return mergedAttributes
}

func writeMetadataJsonFile(fs afero.Afero, path string, attributes map[string]string) error {
	logrus.Infof("Writing metadata attributes into: %s", path)

	attributes = addDeprecatedKeys(attributes)

	jsonData, err := json.Marshal(attributes)
	if err != nil {
		return errors.WithMessage(err, "failed to marshal json")
	}

	return afero.WriteFile(fs, path, jsonData, 0644)
}

func writeMetadataPropertiesFile(fs afero.Afero, path string, attributes map[string]string) error {
	logrus.Infof("Writing metadata attributes into: %s", path)

	attributes = addDeprecatedKeys(attributes)

	// This is done because Go can not guarantee a consistent iteration order for maps, the unit test was impossible to fix as it was too flaky
	sortedKeys := make([]string, 0)
	for key := range attributes {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	var content strings.Builder
	for _, value := range sortedKeys {
		content.WriteString(fmt.Sprintf("%s %s", value, attributes[value]))
		content.WriteString("\n")
	}

	return afero.WriteFile(fs, path, []byte(content.String()), 0644)
}

func addDeprecatedKeys(attributes map[string]string) map[string]string {
	if value, exists := attributes[k8sClusterUidKey]; exists {
		attributes[deprecatedClusterUidKey] = value
	}

	if value, exists := attributes[k8sWorkloadKindKey]; exists {
		attributes[depcrecatedWorkloadKindKey] = value
	}

	if value, exists := attributes[depcrecatedWorkloadNameKey]; exists {
		attributes[depcrecatedWorkloadNameKey] = value
	}

	return attributes
}

func writeContainerConfFile(fs afero.Afero, path string, config map[string]string) error {
	logrus.Infof("Writing container config into: %s", path)

	var containerConfigContent []string

	imageName := generateImageName(config)

	for _, mapping := range attributeMappings {
		value, exists := config[mapping.SourceKey]
		if !exists || mapping.IsImageField {
			continue
		}

		for _, confKey := range mapping.ContainerConfKeys {
			containerConfigContent = append(containerConfigContent, fmt.Sprintf("%s %s", confKey, value))
		}
	}

	if imageName != "" {
		containerConfigContent = append(containerConfigContent, fmt.Sprintf("imageName %s", imageName))
	}

	content := strings.Join(containerConfigContent, "\n") + "\n"

	return afero.WriteFile(fs, path, []byte(content), 0644)
}

func generateImageName(data map[string]string) string {
	var imageName string

	registry, existsRegistry := data[containerImageRegistryKey]
	repository, existsRepository := data[containerImageRepositoryKey]
	tag, existsTag := data[containerImageTagsKey]
	digest, existsDigest := data[containerImageDigestKey]

	if existsRegistry {
		if existsRepository {
			imageName = registry + "/" + repository
		} else {
			imageName = repository
		}

		if existsDigest {
			imageName += "@" + digest
		} else if existsTag {
			imageName += ":" + tag
		}
	}

	return imageName
}
